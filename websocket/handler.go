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
	"github.com/oraichain/orai/x/provider/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"
)

func handleTransaction(c *Context, l *Logger, tx tmtypes.TxResult) {
	l.Debug(":eyes: Inspecting incoming transaction: %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		l.Debug(":alien: Skipping transaction with non-zero code: %d", tx.Result.Code)
		return
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		l.Error(":cold_sweat: Failed to parse transaction logs with error: %s", err.Error())
		return
	}

	for _, log := range logs {
		messageType, err := GetEventValue(log, sdk.EventTypeMessage, sdk.AttributeKeyAction)
		if err != nil {
			l.Error(":cold_sweat: Failed to get message action type with error: %s", err.Error())
			continue
		}

		l.Info(":star: message type: %s", messageType)

		if messageType == (types.MsgSetKYCRequest{}).Type() {
			go handleKYCRequestLog(c, l, log)
		} else if messageType == (types.MsgSetPriceRequest{}).Type() {
			go handlePriceRequestLog(c, l, log)
		} else {
			l.Debug(":ghost: Skipping non-{request/packet} type: %s", messageType)
		} /*else if messageType == (ibc.MsgPacket{}).Type() {
			// Try to get request id from packet. If not then return error.
			_, err := GetEventValue(log, types.EventTypeRequest, types.AttributeKeyID)
			if err != nil {
				l.Debug(":ghost: Skipping non-request packet")
				return
			}
			go handleRequestLog(c, l, log)
		} */
	}
}

func handlePriceRequestLog(c *Context, l *Logger, log sdk.ABCIMessageLog) {
	l.Info(":delivery_truck: Processing incoming request event before checking validators")

	// Skip if not related to this validator
	validators := GetEventValues(log, types.EventTypeSetPriceRequest, types.AttributeRequestValidator)
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

	req, err := GetEventPriceRequest(log)
	if err != nil {
		l.Error(":skull: Failed to parse raw requests with error: %s", err.Error())
	}

	key := <-c.keys
	defer func() {
		c.keys <- key
	}()
	go func(l *Logger, req PriceRequest) {
		// collect data source name from the oScript script
		oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(req.AIRequest.OracleScriptName)

		paths := getPaths(l, oscriptPath)
		aiDataSources := paths[0]
		testCases := paths[1]

		var finalResultStr string
		// create data source results to store in the report
		var dataSourceResultsTest []types.DataSourceResult
		var dataSourceResults []types.DataSourceResult
		var testCaseResults []types.TestCaseResult

		// we have different test cases, so we need to loop through them
		for i := range testCases {
			//put the results from the data sources into the test case to verify if they are good enough
			testCasePath := types.ScriptPath + types.TestCaseStoreKeyString(testCases[i])

			for j := range aiDataSources {
				// Aggregate the required fees for an AI request
				// run the test case script
				cmdTestCase := exec.Command("bash", testCasePath, types.DataSourceStoreKeyString(aiDataSources[j]), req.AIRequest.Input, req.AIRequest.ExpectedOutput)
				var outTestCase bytes.Buffer
				cmdTestCase.Stdout = &outTestCase
				err = cmdTestCase.Run()
				if err != nil {
					l.Error(":skull: failed to execute test case 1st loop: %s", err.Error())
				}

				// collect test case result from the script
				result := strings.TrimSuffix(outTestCase.String(), "\n")

				dataSourceResult := types.NewDataSourceResult(aiDataSources[j], []byte(result), types.ResultSuccess)

				// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
				if result == types.FailedResult || len(result) == 0 {
					// change status to fail so the datasource cannot be rewarded afterwards
					dataSourceResult.Status = types.ResultFailure
				} else {
					// append an data source result into the list
					dataSourceResultsTest = append(dataSourceResultsTest, dataSourceResult)
				}
			}

			// add test case result
			testCaseResults = append(testCaseResults, types.NewTestCaseResult(testCases[i], dataSourceResultsTest))
		}

		// after passing the test cases, we run the actual data sources to collect their results
		// create data source results to store in the report
		// we use dataSourceResultsTest since this list is the complete list of data sources that have passed the test cases
		for i := range dataSourceResultsTest {
			// run the data source script
			cmdTestCase := exec.Command("bash", types.ScriptPath+types.DataSourceStoreKeyString(dataSourceResultsTest[i].Name))
			var outTestCase bytes.Buffer
			cmdTestCase.Stdout = &outTestCase
			err = cmdTestCase.Run()
			if err != nil {
				l.Error(":skull: failed to execute data source script: %s", err.Error())
			}

			// collect test case result from the script
			result := strings.TrimSuffix(outTestCase.String(), "\n")

			fmt.Println("result from data sources: ", result)

			dataSourceResult := types.NewDataSourceResult(dataSourceResultsTest[i].Name, []byte(result), types.ResultSuccess)

			// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
			if len(result) == 0 {
				// change status to fail so the datasource cannot be rewarded afterwards
				dataSourceResult.Status = types.ResultFailure
			}

			// append an data source result into the list
			dataSourceResults = append(dataSourceResults, dataSourceResult)

			// resultInt, _ := strconv.Atoi(result)
			// finalResult = finalResult + resultInt
			finalResultStr = finalResultStr + result + "-"
		}

		fmt.Println("final result string: ", finalResultStr)
		fmt.Println("final result after trimming: ", strings.TrimSuffix(finalResultStr, "-"))
		msgReport := NewReport(req.AIRequest.RequestID, c.validator, dataSourceResults, testCaseResults, key.GetAddress(), sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(5000)))), []byte(finalResultStr))
		if len(finalResultStr) == 0 {
			msgReport.AggregatedResult = []byte(types.FailedResult)
			// Create a new MsgCreateReport to the Oraichain
		} else {
			// "2" here is the expected output that the user wants to get
			cmd := exec.Command("bash", oscriptPath, "aggregation", finalResultStr)
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
		oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(req.AIRequest.OracleScriptName)

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
		var dataSourceResultsTest []types.DataSourceResult
		var finalResult int
		var dataSourceResults []types.DataSourceResult
		var testCaseResults []types.TestCaseResult

		// we have different test cases, so we need to loop through them
		for i := range testCases {
			//put the results from the data sources into the test case to verify if they are good enough
			testCasePath := types.ScriptPath + types.TestCaseStoreKeyString(testCases[i])

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
			cmdTestCase := exec.Command("bash", types.ScriptPath+types.DataSourceStoreKeyString(dataSourceResultsTest[i].Name), imagePath)
			var outTestCase bytes.Buffer
			cmdTestCase.Stdout = &outTestCase
			err = cmdTestCase.Run()
			if err != nil {
				l.Error(":skull: failed to execute data source script: %s", err.Error())
			}

			// collect test case result from the script
			result = strings.TrimSuffix(outTestCase.String(), "\n")

			// append an data source result into the list
			dataSourceResults = append(dataSourceResults, types.NewDataSourceResult(dataSourceResultsTest[i].Name, []byte(result), types.ResultSuccess))

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

func getPaths(l *Logger, oscriptPath string) [][]string {
	//use "data source" as an argument to collect the data source script name
	cmd := exec.Command("bash", oscriptPath, "aiDataSource")
	var dataSourceName bytes.Buffer
	cmd.Stdout = &dataSourceName
	err := cmd.Run()
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

	var listPaths [][]string
	listPaths = append(listPaths, aiDataSources)
	listPaths = append(listPaths, testCases)

	return listPaths
}

// GetEventValues returns the list of all values in the given log with the given type and key.
func GetEventValues(log sdk.ABCIMessageLog, evType string, evKey string) (res []string) {
	for _, ev := range log.Events {
		fmt.Println(":delivery_truck: event collected: ", ev.Type)
		if ev.Type != evType {
			continue
		}

		for _, attr := range ev.Attributes {
			fmt.Println("Attribute key: ", attr.Key)
			fmt.Println("Attribute value: ", attr.Value)
			if attr.Key == evKey {
				res = append(res, attr.Value)
			}
		}
	}
	return res
}

// GetEventValue checks and returns the exact value in the given log with the given type and key.
func GetEventValue(log sdk.ABCIMessageLog, evType string, evKey string) (string, error) {

	values := GetEventValues(log, evType, evKey)
	if len(values) == 0 {
		return "", fmt.Errorf("Cannot find event with type: %s, key: %s", evType, evKey)
	}
	if len(values) > 1 {
		return "", fmt.Errorf("Found more than one event with type: %s, key: %s", evType, evKey)
	}
	return values[0], nil
}

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

// GetEventPriceRequest returns the event price request in the given log.
func GetEventPriceRequest(log sdk.ABCIMessageLog) (PriceRequest, error) {
	requestID, err := GetEventValue(log, types.EventTypeRequestWithData, types.AttributeRequestID)

	req := PriceRequest{}
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

	req = NewPriceRequest(NewAIRequest(requestID, oscriptName, creator, valCount, inputStr, expectedOutputStr))

	fmt.Println("kyc request: ", req)

	return req, nil
}
