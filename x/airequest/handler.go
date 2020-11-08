package airequest

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
	provider "github.com/oraichain/orai/x/provider/exported"
)

// NewHandler creates an sdk.Handler for all the airequest type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		// case types.MsgSetKYCRequest:
		// 	return handleMsgSetKYCRequest(ctx, k, msg)
		case types.MsgSetPriceRequest:
			return handleMsgSetPriceRequest(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgSetPriceRequest(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetPriceRequest) (*sdk.Result, error) {
	validators, err := k.RandomValidators(ctx, msg.MsgAIRequest.ValidatorCount, []byte(msg.MsgAIRequest.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.MsgAIRequest.Fees)
	// Compute the fee allocated for oracle module to distribute to active validators.
	rewardRatio := sdk.NewDecWithPrec(int64(k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)), 2)
	// We need to calculate the final 70% fee given by the user because the remaining 30% must be reserved for the proposer and validators.
	providedCoins, _ := sdk.NewDecCoinsFromCoins(fees...).MulDecTruncate(rewardRatio).TruncateDecimal()

	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.ProviderKeeper.GetDNamesTcNames(msg.MsgAIRequest.OracleScriptName)
	if err != nil {
		return nil, err
	}

	//var finalResult int

	testcaseObjs := make([]provider.TestCaseI, 0)

	dataSourceObjs := make([]provider.AIDataSourceI, 0)

	//var testCaseResults []types.TestCaseResult

	var totalFees sdk.Coins

	// we have different test cases, so we need to loop through them
	for i := 0; i < len(testCases); i++ {
		// loop to run the test case
		// collect all the test cases object to store in the ai request
		testCase, err := k.ProviderKeeper.GetTestCaseI(ctx, testCases[i])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, fmt.Sprintf("failed to get test case: %s", err.Error()))
		}
		testcaseObjs = append(testcaseObjs, testCase)
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(testCase.GetFees()...)
	}

	for j := 0; j < len(aiDataSources); j++ {
		// collect all the data source objects to store in the ai request
		aiDataSource, err := k.ProviderKeeper.GetAIDataSourceI(ctx, aiDataSources[j])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, fmt.Sprintf("failed to get data source: %s", err.Error()))
		}
		dataSourceObjs = append(dataSourceObjs, aiDataSource)
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(aiDataSource.GetFees()...)
	}
	// valFees = 2/5 total dsource and test case fees (70% total in 100% of total fees split into 20% and 50% respectively)
	valFees, _ := sdk.NewDecCoinsFromCoins(totalFees...).MulDec(sdk.NewDecWithPrec(int64(40), 2)).TruncateDecimal()
	// 50% + 20% = 70% * validatorCount fees (since k validators will execute, the fees need to be propotional to the number of vals)
	bothFees := sdk.NewDecCoinsFromCoins(totalFees.Add(valFees...)...)
	finalFees, _ := bothFees.MulDec(sdk.NewDec(int64(msg.MsgAIRequest.ValidatorCount))).TruncateDecimal()
	fmt.Println("both fees: ", bothFees.String())
	fmt.Println("final fees needed: ", finalFees.String())

	// If the total fee is larger than the fee provided by the user then we return error
	if finalFees.IsAnyGT(providedCoins) {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Fees given by the users are less than the total fees needed")
	}
	// set a new request with the aggregated result into blockchain
	request := types.NewAIRequest(msg.MsgAIRequest.RequestID, msg.MsgAIRequest.OracleScriptName, msg.MsgAIRequest.Creator, validators, ctx.BlockHeight(), dataSourceObjs, testcaseObjs, fees, msg.MsgAIRequest.Input, msg.MsgAIRequest.ExpectedOutput)

	fmt.Println("request: ", request)

	//fmt.Printf("request result: %s\n", outOracleScript.String())

	k.SetAIRequest(ctx, request.RequestID, request)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeSetPriceRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
	)
	for _, validator := range validators {
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeRequestValidator, validator.String()),
		)
	}
	ctx.EventManager().EmitEvent(event)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event = sdk.NewEvent(types.EventTypeRequestWithData)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
		sdk.NewAttribute(types.AttributeOracleScriptName, request.OracleScriptName),
		sdk.NewAttribute(types.AttributeRequestCreator, msg.MsgAIRequest.Creator.String()),
		sdk.NewAttribute(types.AttributeRequestValidatorCount, string(msg.MsgAIRequest.ValidatorCount)),
		sdk.NewAttribute(types.AttributeRequestInput, string(msg.MsgAIRequest.Input)),
		sdk.NewAttribute(types.AttributeRequestExpectedOutput, string(msg.MsgAIRequest.ExpectedOutput)),
	)
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// // handleMsgSetKYCRequest is a function message setting a new AI request
// func handleMsgSetKYCRequest(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetKYCRequest) (*sdk.Result, error) {

// 	validators, err := k.RandomValidators(ctx, msg.MsgAIRequest.ValidatorCount, []byte(msg.MsgAIRequest.RequestID))
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
// 	}

// 	// we can safely parse fees to coins since we have validated it in the Msg already
// 	fees, _ := sdk.ParseCoins(msg.MsgAIRequest.Fees)
// 	// Compute the fee allocated for oracle module to distribute to active validators.
// 	rewardRatio := sdk.NewDecWithPrec(int64(k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)), 2)
// 	// We need to calculate the final 70% fee given by the user because the remaining 40% must be reserved for the proposer and validators.
// 	providedCoins, _ := sdk.NewDecCoinsFromCoins(fees...).MulDecTruncate(rewardRatio).TruncateDecimal()

// 	//validatorFees, err := k.GetValidatorFees(ctx, providedCoins, validators)

// 	// collect data source name from the oScript script
// 	oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(msg.MsgAIRequest.OracleScriptName)

// 	//use "data source" as an argument to collect the data source script name
// 	cmd := exec.Command("bash", oscriptPath, "aiDataSource")
// 	cmd.Stdin = strings.NewReader("some input")
// 	var dataSourceName bytes.Buffer
// 	cmd.Stdout = &dataSourceName
// 	err = cmd.Run()
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrFailedToOpenFile, "failed to collect data source name")
// 	}

// 	// collect data source result from the script
// 	result := strings.TrimSuffix(dataSourceName.String(), "\n")

// 	aiDataSources := strings.Fields(result)

// 	//use "test case" as an argument to collect the test case script name
// 	cmd = exec.Command("bash", oscriptPath, "testcase")
// 	cmd.Stdin = strings.NewReader("some input")
// 	var testCaseName bytes.Buffer
// 	cmd.Stdout = &testCaseName
// 	err = cmd.Run()
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(types.ErrFailedToOpenFile, fmt.Sprintf("failed to collect test case name: %s", result))
// 	}

// 	// collect data source result from the script
// 	result = strings.TrimSuffix(testCaseName.String(), "\n")

// 	testCases := strings.Fields(result)

// 	//var finalResult int

// 	var testcaseObjs []TestCase

// 	var dataSourceObjs []AIDataSource

// 	//var testCaseResults []types.TestCaseResult

// 	var totalFees sdk.Coins

// 	// we have different test cases, so we need to loop through them
// 	for i := 0; i < len(testCases); i++ {
// 		// loop to run the test case
// 		// collect all the test cases object to store in the ai request
// 		testCase, err := k.GetTestCase(ctx, testCases[i])
// 		if err != nil {
// 			return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, fmt.Sprintf("failed to get test case: %s", err.Error()))
// 		}
// 		testcaseObjs = append(testcaseObjs, testCase)
// 		// Aggregate the required fees for an AI request
// 		totalFees = totalFees.Add(testCase.Fees...)
// 	}

// 	for j := 0; j < len(aiDataSources); j++ {
// 		// collect all the data source objects to store in the ai request
// 		aiDataSource, err := k.GetAIDataSource(ctx, aiDataSources[j])
// 		if err != nil {
// 			return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, fmt.Sprintf("failed to get data source: %s", err.Error()))
// 		}
// 		dataSourceObjs = append(dataSourceObjs, aiDataSource)
// 		// Aggregate the required fees for an AI request
// 		totalFees = totalFees.Add(aiDataSource.Fees...)
// 	}

// 	// If the total fee is larger than the fee provided by the user then we return error
// 	if totalFees.IsAnyGT(providedCoins) {
// 		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Fees given by the users are less than the total fees needed")
// 	}

// 	// set a new request with the aggregated result into blockchain
// 	request := types.NewAIRequest(msg.MsgAIRequest.RequestID, msg.MsgAIRequest.OracleScriptName, msg.MsgAIRequest.Creator, validators, ctx.BlockHeight(), dataSourceObjs, testcaseObjs, fees, msg.MsgAIRequest.Input, msg.MsgAIRequest.ExpectedOutput)

// 	//fmt.Printf("request result: %s\n", outOracleScript.String())

// 	k.SetAIRequest(ctx, request.RequestID, request)

// 	// TODO: Define your msg events
// 	// Emit an event describing a data request and asked validators.
// 	event := sdk.NewEvent(types.EventTypeSetKYCRequest)
// 	event = event.AppendAttributes(
// 		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
// 	)
// 	for _, validator := range validators {
// 		event = event.AppendAttributes(
// 			sdk.NewAttribute(types.AttributeRequestValidator, validator.String()),
// 		)
// 	}
// 	ctx.EventManager().EmitEvent(event)

// 	// TODO: Define your msg events
// 	// Emit an event describing a data request and asked validators.
// 	event = sdk.NewEvent(types.EventTypeRequestWithData)
// 	event = event.AppendAttributes(
// 		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
// 		sdk.NewAttribute(types.AttributeOracleScriptName, request.OracleScriptName),
// 		sdk.NewAttribute(types.AttributeRequestCreator, msg.MsgAIRequest.Creator.String()),
// 		sdk.NewAttribute(types.AttributeRequestImageHash, msg.ImageHash),
// 		sdk.NewAttribute(types.AttributeRequestImageName, msg.ImageName),
// 		sdk.NewAttribute(types.AttributeRequestValidatorCount, string(msg.MsgAIRequest.ValidatorCount)),
// 		sdk.NewAttribute(types.AttributeRequestInput, msg.MsgAIRequest.Input),
// 		sdk.NewAttribute(types.AttributeRequestExpectedOutput, msg.MsgAIRequest.ExpectedOutput),
// 	)
// 	ctx.EventManager().EmitEvent(event)
// 	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
// }
