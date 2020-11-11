package websocket

import (
	"bytes"
	"fmt"
	"io"
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

// GetEventKYCRequest returns the event kyc request in the given log.
func GetEventKYCRequest(log sdk.ABCIMessageLog) (KYCRequest, error) {
	requestID, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestID)
	req := KYCRequest{}
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

	req = NewKYCRequest(imageHash, imageName, NewAIRequest(requestID, oscriptName, creator, valCount, inputStr, expectedOutputStr))

	fmt.Println("kyc request: ", req)

	return req, nil
}

func handleKYCRequestLog(c *Context, l *Logger, log sdk.ABCIMessageLog) {

	l.Info(":delivery_truck: Processing incoming request event before checking validators")

	// Skip if not related to this validator
	validators := GetEventValues(log, types.EventTypeSetKYCRequest, types.AttributeRequestValidator)
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

	req, err := GetEventKYCRequest(log)
	if err != nil {
		l.Error(":skull: Failed to parse raw requests with error: %s", err.Error())
	}

	key := <-c.keys
	defer func() {
		c.keys <- key
	}()
	go func(l *Logger, req KYCRequest) {
		// collect the image from IPFS for all nodes that handle this req
		imagePath := ".images/" + req.ImageName

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
		var finalResult int
		var dataSourceResults []exported.DataSourceResultI
		var testCaseResults []exported.TestCaseResultI

		// we have different test cases, so we need to loop through them
		for i := range testCases {
			//put the results from the data sources into the test case to verify if they are good enough
			testCasePath := provider.ScriptPath + provider.TestCaseStoreKeyString(testCases[i])

			for j := range aiDataSources {
				// Aggregate the required fees for an AI request
				// run the test case script
				cmdTestCase := exec.Command("bash", testCasePath, aiDataSources[j], imagePath, string(req.AIRequest.ExpectedOutput))
				var outTestCase bytes.Buffer
				cmdTestCase.Stdout = &outTestCase
				err = cmdTestCase.Run()
				if err != nil {
					l.Error(":skull: failed to execute test case 1st loop: %s", err.Error())
				}

				// collect test case result from the script
				result = strings.TrimSuffix(outTestCase.String(), "\n")

				dataSourceResult := types.NewDataSourceResult(aiDataSources[j], []byte(result), types.ResultSuccess)

				// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
				if result == types.FailedResult || len(result) == 0 {
					// change status to fail so the datasource cannot be rewarded afterwards
					dataSourceResult.Status = types.ResultFailure
				}

				// append an data source result into the list
				dataSourceResultsTest = append(dataSourceResultsTest, dataSourceResult)
			}

			// add test case result
			testCaseResults = append(testCaseResults, types.NewTestCaseResult(testCases[i], dataSourceResultsTest))
		}

		// after passing the test cases, we run the actual data sources to collect their results
		// create data source results to store in the report
		for i := range dataSourceResultsTest {
			// run the data source script
			cmdTestCase := exec.Command("bash", provider.ScriptPath+provider.DataSourceStoreKeyString(dataSourceResultsTest[i].GetName()), imagePath)
			var outTestCase bytes.Buffer
			cmdTestCase.Stdout = &outTestCase
			err = cmdTestCase.Run()
			if err != nil {
				l.Error(":skull: failed to execute data source script: %s", err.Error())
			}

			// collect test case result from the script
			result = strings.TrimSuffix(outTestCase.String(), "\n")

			// append an data source result into the list
			dataSourceResults = append(dataSourceResults, types.NewDataSourceResult(dataSourceResultsTest[i].GetName(), []byte(result), types.ResultSuccess))

			resultInt, _ := strconv.Atoi(result)
			finalResult = finalResult + resultInt
		}

		// report := k.GetAllReports(ctx)
		// fmt.Printf("Report: %v\n", report)

		// "2" here is the expected output that the user wants to get
		cmd = exec.Command("bash", oscriptPath, "aggregation", strconv.Itoa(finalResult), "2")
		var res bytes.Buffer
		cmd.Stdout = &res
		err = cmd.Run()
		if err != nil {
			l.Error(":skull: failed to aggregate results: %s", err.Error())
		}

		// collect data source result from the script
		ress := strings.TrimSuffix(res.String(), "\n")
		fmt.Printf("final result from oScript: %s\n", ress)

		msgReport := NewReport(req.AIRequest.RequestID, c.validator, dataSourceResults, testCaseResults, key.GetAddress(), sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(5000)))), []byte(ress))

		// Create a new MsgCreateReport to the Oraichain
		SubmitReport(c, l, key, msgReport)

	}(l.With("reqid", req.AIRequest.RequestID, "oscriptname", req.AIRequest.OracleScriptName), req)
}
