package subscribe

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

func (subscriber *Subscriber) checkRequestFees(aiOracle *types.AIOracle, queryClient types.QueryClient) (sdk.Coins, bool) {
	if !aiOracle.Fees.IsZero() && !aiOracle.Fees.Empty() {
		oneValFees := sdk.NewDecCoinsFromCoins(aiOracle.Fees...).QuoDec(sdk.NewDec(int64(len(aiOracle.GetValidators()))))
		resp, _ := queryClient.Param(context.Background(), &types.QueryParamRequest{Param: "KeyAIOracleRewardPercentages"})
		aiOracleRewardPercentage := resp.Param
		aiOracleRewardPercentage += 100
		providerFees := oneValFees.QuoDec(sdk.NewDecWithPrec(int64(aiOracleRewardPercentage), 2))
		valFees, _ := providerFees.MulDec(sdk.NewDecWithPrec(int64(resp.Param), 2)).TruncateDecimal()
		if valFees.IsAllLT(subscriber.config.RequestFees) {
			return valFees, false
		}
		return valFees, true
	}
	return sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(0))), false
}

func (subscriber *Subscriber) handleAIRequestLog(queryClient types.QueryClient, attrs []sdk.Attribute) ([]byte, error) {
	subscriber.log.Info(":delivery_truck: Processing incoming request event before checking validators")

	aiOracle := &types.AIOracle{}
	reportTypes := 0
	for _, attr := range attrs {
		switch attr.Key {
		case types.AttributeRequest:
			err := json.Unmarshal([]byte(attr.Value), aiOracle)
			if err != nil {
				subscriber.log.Error(":skull: Cannot decode request to execute: %s", attr.Value)
				return nil, err
			}
			break
		case types.AttributeReport:
			reportTypes = 0
			break
		case types.AttributeBaseReport:
			reportTypes = 1
		default:
			reportTypes = -1
			break
		}
	}
	fmt.Println("ai oracle data: ", aiOracle)
	// Skip if not related to this validator
	subscriber.log.Info(":delivery_truck: Validator lists: ", aiOracle.GetValidators())
	currentValidator := sdk.ValAddress(subscriber.cliCtx.GetFromAddress())
	hasMe := false
	for _, validator := range aiOracle.GetValidators() {
		subscriber.log.Debug(":delivery_truck: validator: %s", validator)
		if validator.Equals(currentValidator) {
			hasMe = true
			break
		}
	}
	valFees, isSufficient := subscriber.checkRequestFees(aiOracle, queryClient)
	if !isSufficient {
		subscriber.log.Info(":delivery_truck: request pays you: %v less fees than required which is: %v. Stop executing the request", valFees, subscriber.config.RequestFees)
		return nil, fmt.Errorf("request pays less fees than required. Stop executing the request")
	}
	ctx := context.Background()
	var reportBytes []byte
	var err error
	if hasMe {
		switch reportTypes {
		case 0:
			reportBytes, err = subscriber.ExecuteAIOracle(queryClient, ctx, *aiOracle, currentValidator)
			break
		case 1:
			reportBytes, err = subscriber.ExecuteAIOracleTestOnly(queryClient, ctx, *aiOracle, currentValidator)
			break
		default:
			return nil, fmt.Errorf("Invalid report types")
		}
	}
	subscriber.log.Debug(":delivery_truck: not your request to execute. Your address is: %v, where the list of validators are: %v", currentValidator, aiOracle.GetValidators())
	return reportBytes, err
}

func (subscriber *Subscriber) ExecuteAIOracle(queryClient types.QueryClient, ctx context.Context, aiOracle types.AIOracle, valAddress sdk.ValAddress) ([]byte, error) {
	// querier to interact with the wasm contract
	var dataSourceResults []*types.Result
	var resultArr = []string{}
	// collect list entries to get entry length
	entries, err := queryClient.DataSourceEntries(ctx, &types.QueryDataSourceEntriesContract{
		Contract: aiOracle.GetContract(),
		Request:  &types.EmptyParams{},
	})
	if err != nil {
		subscriber.log.Error(fmt.Sprintf("Cannot get data source entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
		return nil, err
	}

	// if there's no data source we stop executing and return no report
	if len(entries.GetData()) <= 0 {
		subscriber.log.Error(fmt.Sprintf("The data source entry list is empty"))
		return nil, err
	}

	// loop to execute data source one by one
	for _, entry := range entries.GetData() {
		// run the data source script
		subscriber.log.Info(fmt.Sprintf("Data source entrypoint: %v and input: %v", entry, string(aiOracle.GetInput())))
		result, err := queryClient.DataSourceContract(ctx, &types.QueryDataSourceContract{
			Contract: aiOracle.GetContract(),
			Request: &types.RequestDataSource{
				Dsource: entry,
				Input:   string(aiOracle.GetInput()),
			},
		})
		// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
		dataSourceResult := types.NewResult(entry, result.GetData(), types.ResultSuccess)
		if err != nil {
			subscriber.log.Error(fmt.Sprintf("Cannot execute given data source %v with error: %v", entry.Url, err))
			// change status to fail so the datasource cannot be rewarded afterwards
			dataSourceResult.Status = types.ResultFailure
			dataSourceResult.Result = []byte(types.FailedResponseDs)
			continue
		}
		resultArr = append(resultArr, string(result.Data))
		// append an data source result into the list
		dataSourceResults = append(dataSourceResults, dataSourceResult)
	}

	// store report into blockchain as proof for executing AI requests
	report := types.NewReport(aiOracle.GetRequestID(), dataSourceResults, aiOracle.GetBlockHeight(), []byte{}, valAddress, types.ResultSuccess, sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(0))))

	subscriber.log.Info(fmt.Sprintf("results collected from the data sources: %v", resultArr))
	aggregatedResult, err := queryClient.OracleScriptContract(ctx, &types.QueryOracleScriptContract{
		Contract: aiOracle.GetContract(),
		Request: &types.RequestOracleScript{
			Results: resultArr,
		},
	})
	if err != nil {
		subscriber.log.Error(fmt.Sprintf("Cannot execute oracle script contract %v with error: %v", aiOracle.GetContract(), err))
		report.AggregatedResult = []byte(types.FailedResponseOs)
		report.ResultStatus = types.ResultFailure
	}
	report.AggregatedResult = aggregatedResult.Data
	subscriber.log.Info(fmt.Sprintf("Oracle script final result: %v", aggregatedResult))
	reportBytes, err := json.Marshal(report)
	if err != nil {
		subscriber.log.Error(fmt.Sprintf("Cannot marshal report %v with error: %v", report, err))
		return nil, err
	}
	subscriber.log.Info(fmt.Sprintf("Finish handling the AI oracles with report: %v", report))
	return reportBytes, err
}

func (subscriber *Subscriber) ExecuteAIOracleTestOnly(queryClient types.QueryClient, ctx context.Context, aiOracle types.AIOracle, valAddress sdk.ValAddress) ([]byte, error) {
	// querier to interact with the wasm contract
	var resultsWithTc []*types.ResultWithTestCase
	// collect list entries to get entry length
	entries, err := queryClient.DataSourceEntries(ctx, &types.QueryDataSourceEntriesContract{
		Contract: aiOracle.GetContract(),
		Request:  &types.EmptyParams{},
	})
	if err != nil {
		subscriber.log.Error(fmt.Sprintf("Cannot get data source entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
		return nil, err
	}
	tCaseEntries, err := queryClient.TestCaseEntries(ctx, &types.QueryTestCaseEntriesContract{
		Contract: aiOracle.GetContract(),
		Request:  &types.EmptyParams{},
	})
	if err != nil {
		subscriber.log.Error(fmt.Sprintf("Cannot get test case entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
		return nil, err
	}

	// if there's no data source or test case then we stop executing
	if len(entries.GetData()) <= 0 || len(tCaseEntries.GetData()) <= 0 {
		subscriber.log.Error(fmt.Sprintf("The data source or test case entry list is empty"))
		return nil, err
	}

	// loop to execute data source one by one
	for _, entry := range entries.GetData() {
		var results []*types.Result
		for _, tCaseEntry := range tCaseEntries.GetData() {
			// run the data source script
			subscriber.log.Info(fmt.Sprintf("Data source entrypoint: %v and input: %v", entry, string(aiOracle.GetInput())))
			subscriber.log.Info(fmt.Sprintf("Testcase entrypoint: %v", tCaseEntry))
			result, err := queryClient.TestCaseContract(ctx, &types.QueryTestCaseContract{
				Contract: aiOracle.GetContract(),
				Request: &types.RequestTestCase{
					Tcase: tCaseEntry,
					Input: entry,
				},
			})
			tCaseResult := types.NewResult(tCaseEntry, []byte{}, types.ResultSuccess)
			if err != nil {
				subscriber.log.Error(fmt.Sprintf("Cannot execute test case %v with error: %v", tCaseEntry.Url, err))
				tCaseResult.Result = []byte(types.FailedResponseTc)
				tCaseResult.Status = types.ResultFailure
				results = append(results, tCaseResult)
				continue
			}
			tCaseResult.Result = result.GetData()
			results = append(results, tCaseResult)
		}
		resultWithTc := types.NewResultWithTestCase(entry, results, types.ResultSuccess)
		resultsWithTc = append(resultsWithTc, resultWithTc)
	}
	// store report into blockchain as proof for executing AI requests
	report := types.NewTestCaseReport(aiOracle.GetRequestID(), resultsWithTc, aiOracle.GetBlockHeight(), valAddress, sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(0))))
	reportBytes, err := json.Marshal(report)
	if err != nil {
		subscriber.log.Error(fmt.Sprintf("Cannot marshal report %v with error: %v", report, err))
	}
	return reportBytes, nil
}
