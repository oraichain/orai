package subscribe

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	aiRequest "github.com/oraichain/orai/x/airequest"

	"time"

	"github.com/cosmos/cosmos-sdk/client"
	utils "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/oraichain/orai/x/websocket/types"
)

// GetEventAIRequest returns the event AI request in the given log.
func GetEventAIRequest(log sdk.ABCIMessageLog) (*aiRequest.AIRequest, error) {
	// ev := aiRequest.EventTypeRequestWithData
	// requestID, err := GetEventValue(log, ev, aiRequest.AttributeRequestID)

	// if err != nil {
	// 	return nil, err
	// }
	// oscriptName, err := GetEventValue(log, ev, aiRequest.AttributeOracleScriptName)
	// if err != nil {
	// 	return nil, err
	// }
	// creatorStr, err := GetEventValue(log, ev, aiRequest.AttributeRequestCreator)
	// if err != nil {
	// 	return nil, err
	// }

	// valCountStr, err := GetEventValue(log, ev, aiRequest.AttributeRequestValidatorCount)
	// if err != nil {
	// 	return nil, err
	// }

	// inputStr, err := GetEventValue(log, ev, aiRequest.AttributeRequestInput)

	// expectedOutputStr, err := GetEventValue(log, ev, aiRequest.AttributeRequestExpectedOutput)

	// creator, err := sdk.AccAddressFromBech32(creatorStr)
	// if err != nil {
	// 	return nil, err
	// }

	// valCount, _ := strconv.ParseInt(valCountStr, 10, 64)

	// req := aiRequest.NewAIRequest(requestID, oscriptName, creator, valCount, inputStr, expectedOutputStr)

	// return req, nil
	return nil, nil
}

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func (subscriber *Subscriber) submitReport(cliCtx *client.Context, config *types.WebSocketConfig, msgReport *types.MsgCreateReport) {

	txHash := ""
	if err := msgReport.ValidateBasic(); err != nil {
		subscriber.log.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return
	}

	for try := uint64(1); try <= config.MaxTry; try++ {
		subscriber.log.Info(":e-mail: Try to broadcast report transaction(%d/%d)", try, config.MaxTry)

		txBldr := cliCtx.TxConfig.NewTxBuilder()
		txBldr.SetMemo("websocket")

		bytes, err := cliCtx.JSONMarshaler.MarshalJSON(msgReport)

		res, err := cliCtx.BroadcastTxSync(bytes)
		subscriber.log.Info(":star: response: %v", res)
		if err == nil {
			txHash = res.TxHash
			break
		}
		subscriber.log.Info(":warning: Failed to broadcast tx with error: %s", err.Error())
		time.Sleep(config.RPCPollInterval)
	}
	if txHash == "" {
		subscriber.log.Error(":exploding_head: Cannot try to broadcast more than %d try", config.MaxTry)
		return
	}
	for start := time.Now(); time.Since(start) < config.BroadcastTimeout; {
		time.Sleep(config.RPCPollInterval)
		txRes, err := utils.QueryTx(*cliCtx, txHash)
		if err != nil {
			subscriber.log.Debug(":warning: Failed to query tx with error: %s", err.Error())
			continue
		}
		if txRes.Code != 0 {
			subscriber.log.Error(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", txRes.Code, txRes.RawLog, txRes.TxHash)
			return
		}
		subscriber.log.Info(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", txHash)
		return
	}
	subscriber.log.Info(":question_mark: Cannot get transaction response from hash: %s transaction might be included in the next few blocks or check your node's health.", txHash)
}

func (subscriber *Subscriber) handleAIRequestLog(cliCtx *client.Context, queryClient types.QueryClient, attrMap map[string]string) {
	fmt.Println(":delivery_truck: Processing incoming request event before checking validators")

	// // Skip if not related to this validator
	// validators := GetEventValues(log, aiRequest.EventTypeSetAIRequest, aiRequest.AttributeRequestValidator)
	// hasMe := false
	// for _, validator := range validators {
	// 	fmt.Println(":delivery_truck: validator: %s", validator)
	// 	if validator == c.validator.String() {
	// 		hasMe = true
	// 		break
	// 	}
	// }

	// if !hasMe {
	// 	fmt.Println(":next_track_button: Skip request not related to this validator")
	// 	return
	// }

	fmt.Println(":delivery_truck: Processing incoming request event")

	// req, err := GetEventAIRequest(log)
	// if err != nil {
	// 	fmt.Println(":skull: Failed to parse raw requests with error: %s", err.Error())
	// }

	// key := <-c.keys
	// defer func() {
	// 	c.keys <- key
	// }()

	// // invoke a new goroutine to run the thread in parallel
	// // collect data source name from the oScript script

	// // encode the input and output back to base64 type to forward to the test case
	// //input := base64.StdEncoding.EncodeToString([]byte(req.Input))
	// //expectedOutput := base64.StdEncoding.EncodeToString([]byte(req.ExpectedOutput))

	// // collect ai data sources and test cases from the ai request event.
	// aiDataSources, testCases, err := getPaths(log)

	// var finalResultStr string
	// // create data source results to store in the report
	// var dataSourceResultsTest []*websocket.DataSourceResult
	// var dataSourceResults []*websocket.DataSourceResult
	// var testCaseResults []*websocket.TestCaseResult

	// // we have different test cases, so we need to loop through them
	// for i := range testCases {
	// 	//put the results from the data sources into the test case to verify if they are good enough
	// 	for j := range aiDataSources {
	// 		//// collect test case result from the script
	// 		// outTestCase, err := ExecPythonFile("python", getTCasePath(testCases[i]), []string{provider.DataSourceStoreKeyString(aiDataSources[j]), req.Input, req.ExpectedOutput})
	// 		result := "from wasm keeper"
	// 		if err != nil {
	// 			fmt.Println(":skull: failed to execute test case 1st loop: %s", err.Error())
	// 		}

	// 		fmt.Println("result after running test case: ", result)

	// 		//fmt.Info("star: result after running test case: ", result)

	// 		dataSourceResult := websocket.NewDataSourceResult(aiDataSources[j], []byte(result), websocket.ResultSuccess)

	// 		// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
	// 		if result == websocket.FailedResult || len(result) == 0 {
	// 			// change status to fail so the datasource cannot be rewarded afterwards
	// 			dataSourceResult.Status = websocket.ResultFailure
	// 			dataSourceResult.Result = []byte(websocket.FailedResponseTc)
	// 		}
	// 		// append an data source result into the list
	// 		dataSourceResultsTest = append(dataSourceResultsTest, dataSourceResult)
	// 	}

	// 	// add test case result
	// 	testCaseResults = append(testCaseResults, websocket.NewTestCaseResult(testCases[i], dataSourceResultsTest))
	// }

	// // after passing the test cases, we run the actual data sources to collect their results
	// // create data source results to store in the report
	// // we use dataSourceResultsTest since this list is the complete list of data sources that have passed the test cases
	// for i := range dataSourceResultsTest {
	// 	// run the data source script

	// 	var dataSourceResult *websocket.DataSourceResult
	// 	if dataSourceResultsTest[i].GetStatus() == websocket.ResultSuccess {
	// 		// outDataSource, err = ExecPythonFile("python", getDSourcePath(dataSourceResultsTest[i].GetName()), []string{})
	// 		result := "excute data source contract"

	// 		//result = strings.TrimSuffix(result, "\r")
	// 		fmt.Println("star: result from data sources: ", result)
	// 		// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
	// 		dataSourceResult = websocket.NewDataSourceResult(dataSourceResultsTest[i].GetName(), []byte(result), websocket.ResultSuccess)
	// 		if result == websocket.FailedResult || len(result) == 0 {
	// 			// change status to fail so the datasource cannot be rewarded afterwards
	// 			dataSourceResult.Status = websocket.ResultFailure
	// 			dataSourceResult.Result = []byte(websocket.FailedResponseDs)
	// 		} else {
	// 			//resultArr = append(resultArr, resultFloat)
	// 			finalResultStr = finalResultStr + result
	// 		}
	// 	} else {
	// 		dataSourceResult = websocket.NewDataSourceResult(dataSourceResultsTest[i].GetName(), []byte(dataSourceResultsTest[i].GetResult()), types.ResultFailure)
	// 	}
	// 	// append an data source result into the list
	// 	dataSourceResults = append(dataSourceResults, dataSourceResult)
	// }
	// finalResultStr = strings.TrimSuffix(finalResultStr, "-")
	// fmt.Println("star: final result after trimming: ", finalResultStr)
	// // Create a new MsgCreateReport with a new reporter to the Oraichain
	// reporter := websocket.NewReporter(key.GetAddress(), key.GetName(), c.validator)
	// msgReport := websocket.NewMsgCreateReport(req.RequestID, dataSourceResults, testCaseResults, reporter, sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(5000)))), []byte(finalResultStr), types.ResultSuccess)
	// if len(finalResultStr) == 0 {
	// 	msgReport.AggregatedResult = []byte(websocket.FailedResponseOs)
	// 	msgReport.ResultStatus = websocket.ResultFailure
	// 	// Create a new MsgCreateReport to the Oraichain
	// } else {
	// 	ress := "exec oracle contract"
	// 	fmt.Printf("final result from oScript: %s\n", ress)
	// 	msgReport.AggregatedResult = []byte(ress)
	// }

	// // TODO: check aggregated result value of the oracle script to verify if it fails or success
	// subscriber.submitReport(c, key, msgReport)

}
