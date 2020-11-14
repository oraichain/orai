package websocket

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	provider "github.com/oraichain/orai/x/provider/types"
	"github.com/oraichain/orai/x/websocket/exported"
	"github.com/oraichain/orai/x/websocket/types"
)

// GetEventOCRRequest returns the event ocr request in the given log.
func GetEventOCRRequest(log sdk.ABCIMessageLog) (OCRRequest, error) {
	requestID, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestID)
	req := OCRRequest{}
	if err != nil {
		return req, err
	}
	oscriptName, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeOracleScriptName)
	if err != nil {
		return req, err
	}
	creatorStr, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestCreator)
	if err != nil {
		return req, err
	}
	imageHash, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestImageHash)
	if err != nil {
		return req, err
	}
	imageName, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestImageName)
	if err != nil {
		return req, err
	}
	valCountStr, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestValidatorCount)
	if err != nil {
		return req, err
	}

	inputStr, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestInput)

	expectedOutputStr, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestExpectedOutput)

	fmt.Println("expected output str: ", expectedOutputStr)

	creator, err := sdk.AccAddressFromBech32(creatorStr)
	if err != nil {
		return req, err
	}

	valCount, _ := strconv.ParseInt(valCountStr, 10, 64)

	req = NewOCRRequest(imageHash, imageName, NewAIRequest(requestID, oscriptName, creator, valCount, inputStr, expectedOutputStr))

	fmt.Println("classification request: ", req)

	return req, nil
}

func handleOCRRequestLog(c *Context, l *Logger, log sdk.ABCIMessageLog) {

	l.Info(":delivery_truck: Processing incoming request event before checking validators")

	// Skip if not related to this validator
	validators := GetEventValues(log, types.EventTypeSetOCRRequest, types.AttributeRequestValidator)
	hasMe := false
	for _, validator := range validators {
		l.Info(":delivery_truck: validator: %s", validator)
		if validator == c.validator.String() {
			hasMe = true
			break
		}
	}

	if !hasMe {
		l.Debug(":next_track_button: Skip request not related to this validator")
		return
	}

	l.Info(":delivery_truck: Processing incoming request event")

	req, err := GetEventClassificationRequest(log)
	if err != nil {
		l.Error(":skull: Failed to parse raw requests with error: %s", err.Error())
	}

	key := <-c.keys
	defer func() {
		c.keys <- key
	}()
	go func(l *Logger, req ClassificationRequest) {

		dir, err := ioutil.TempDir("./", ".images")
		if err != nil {
			l.Error(":skull: Failed to create directory from image path: %s", err.Error())
		}
		defer os.RemoveAll(dir)

		// collect the image from IPFS for all nodes that handle this req
		imagePath := dir + "/" + req.ImageName
		fmt.Println("image path: ", imagePath)

		out, err := os.Create(imagePath)
		if err != nil {
			l.Error(":skull: Failed to create output from image path: %s", err.Error())
		}
		defer out.Close()
		resp, err := http.Post(types.IPFSUrl+types.IPFSCat+"?arg="+req.ImageHash, "application/json", nil)
		if err != nil {
			l.Error(":skull: Failed to receive response from IPFS: %s", err.Error())
		}
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)

		if err != nil {
			l.Error(":skull: Failed to create a new image from IPFS: %s", err.Error())
		}

		// collect data source name from the oScript script
		oscriptPath := provider.ScriptPath + provider.OracleScriptStoreKeyString(req.AIRequest.OracleScriptName)

		//use "data source" as an argument to collect the data source script name
		cmd := exec.Command("bash", oscriptPath, "aiDataSource")
		var dataSourceName bytes.Buffer
		cmd.Stdout = &dataSourceName
		err = cmd.Run()
		if err != nil {
			l.Error(":skull: failed to collect data source name: %s", err.Error())
		}

		// collect data source result from the script
		result := strings.TrimSuffix(dataSourceName.String(), "\n")

		aiDataSources := strings.Fields(result)

		//use "test case" as an argument to collect the test case script name
		cmd = exec.Command("bash", oscriptPath, "testcase")
		var testCaseName bytes.Buffer
		cmd.Stdout = &testCaseName
		err = cmd.Run()
		if err != nil {
			l.Error(":skull: failed to collect test case name: %s", err.Error())
		}

		// collect data source result from the script
		result = strings.TrimSuffix(testCaseName.String(), "\n")

		testCases := strings.Fields(result)

		// create data source results to store in the report
		var dataSourceResultsTest []exported.DataSourceResultI
		var finalResultStr string
		var dataSourceResults []exported.DataSourceResultI
		var testCaseResults []exported.TestCaseResultI

		// we have different test cases, so we need to loop through them
		for i := range testCases {
			//put the results from the data sources into the test case to verify if they are good enough
			for j := range aiDataSources {
				// Aggregate the required fees for an AI request
				// run the test case script
				fmt.Println("test case path: ", getTCasePath(testCases[i])+provider.DataSourceStoreKeyString(aiDataSources[j]))
				cmdTestCase := exec.Command("bash", getTCasePath(testCases[i]), provider.DataSourceStoreKeyString(aiDataSources[j]), imagePath, req.AIRequest.ExpectedOutput, req.AIRequest.Input)
				var outTestCase bytes.Buffer
				cmdTestCase.Stdout = &outTestCase
				err = cmdTestCase.Run()
				if err != nil {
					l.Error(":skull: failed to execute test case 1st loop: %s", err.Error())
				}

				// collect test case result from the script
				result := strings.TrimSuffix(outTestCase.String(), "\n")

				fmt.Println("result after running test case: ", result)

				dataSourceResult := types.NewDataSourceResult(aiDataSources[j], []byte(result), types.ResultSuccess)

				// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
				if result == types.FailedResult || len(result) == 0 {
					// change status to fail so the datasource cannot be rewarded afterwards
					dataSourceResult.Status = types.ResultFailure
					dataSourceResult.Result = []byte(types.FailedResponseTc)
				}
				// append an data source result into the list
				dataSourceResultsTest = append(dataSourceResultsTest, dataSourceResult)
			}

			// add test case result
			testCaseResults = append(testCaseResults, types.NewTestCaseResult(testCases[i], dataSourceResultsTest))
		}

		// after passing the test cases, we run the actual data sources to collect their results
		// create data source results to store in the report
		// we use dataSourceResultsTest since this list is the complete list of data sources that have passed the test cases
		dataSourceCount := 0
		for i := range dataSourceResultsTest {
			// run the data source script
			var outTestCase bytes.Buffer
			var dataSourceResult types.DataSourceResult
			if dataSourceResultsTest[i].GetStatus() == types.ResultSuccess {
				cmdTestCase := exec.Command("bash", getDSourcePath(dataSourceResultsTest[i].GetName()), imagePath, req.AIRequest.Input)
				cmdTestCase.Stdout = &outTestCase
				err = cmdTestCase.Run()
				if err != nil {
					l.Error(":skull: failed to execute data source script: %s", err.Error())
				}
				// collect test case result from the script
				result := strings.TrimSuffix(outTestCase.String(), "\n")
				fmt.Println("result from data sources: ", result)
				// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
				dataSourceResult = types.NewDataSourceResult(dataSourceResultsTest[i].GetName(), []byte(result), types.ResultSuccess)
				if result == types.FailedResult || len(result) == 0 {
					// change status to fail so the datasource cannot be rewarded afterwards
					dataSourceResult.Status = types.ResultFailure
					dataSourceResult.Result = []byte(types.FailedResponseDs)
				} else {
					if dataSourceCount == 0 {
						finalResultStr += result
					}
					dataSourceCount++
				}
			} else {
				dataSourceResult = types.NewDataSourceResult(dataSourceResultsTest[i].GetName(), []byte(dataSourceResultsTest[i].GetResult()), types.ResultFailure)
			}
			// append an data source result into the list
			dataSourceResults = append(dataSourceResults, dataSourceResult)
		}
		msgReport := NewReport(req.AIRequest.RequestID, c.validator, dataSourceResults, testCaseResults, key.GetAddress(), sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(5000)))), []byte(finalResultStr))
		if len(finalResultStr) == 0 {
			msgReport.AggregatedResult = []byte(types.FailedResponseOs)
			// Create a new MsgCreateReport to the Oraichain
		} else {
			// "2" here is the expected output that the user wants to get
			cmd := exec.Command("bash", oscriptPath, "aggregation", finalResultStr, fmt.Sprint(dataSourceCount))
			var res bytes.Buffer
			cmd.Stdout = &res
			err = cmd.Run()
			if err != nil {
				l.Error(":skull: failed to aggregate results: %s", err.Error())
			}

			// collect data source result from the script
			ress := strings.TrimSuffix(res.String(), "\n")
			fmt.Printf("final result from oScript: %s\n", ress)
			msgReport.AggregatedResult = []byte(ress)
		}

		// Create a new MsgCreateReport to the Oraichain
		SubmitReport(c, l, key, msgReport)

	}(l.With("reqid", req.AIRequest.RequestID, "oscriptname", req.AIRequest.OracleScriptName), req)
}
