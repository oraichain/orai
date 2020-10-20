package provider

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/provider/keeper"
	"github.com/oraichain/orai/x/provider/types"
)

// NewHandler creates an sdk.Handler for all the provider type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateOracleScript:
			return handleMsgCreateOracleScript(ctx, k, msg)
		case types.MsgCreateAIDataSource:
			return handleMsgCreateAIDataSource(ctx, k, msg)
		case types.MsgEditOracleScript:
			return handleMsgEditOracleScript(ctx, k, msg)
		case types.MsgEditAIDataSource:
			return handleMsgEditAIDataSource(ctx, k, msg)
		case types.MsgSetKYCRequest:
			return handleMsgSetKYCRequest(ctx, k, msg)
		case types.MsgSetPriceRequest:
			return handleMsgSetPriceRequest(ctx, k, msg)
		case types.MsgCreateTestCase:
			return handleMsgCreateTestCase(ctx, k, msg)
		case types.MsgEditTestCase:
			return handleMsgEditTestCase(ctx, k, msg)
		case types.MsgCreateReport:
			return handleMsgAddReport(ctx, k, msg)
		case types.MsgAddReporter:
			return handleMsgAddReporter(ctx, k, msg)
		case types.MsgRemoveReporter:
			return handleMsgRemoveReporter(ctx, k, msg)
		case types.MsgCreateStrategy:
			return handleMsgCreateStrategy(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateOracleScript is a function message setting oScript
func handleMsgCreateOracleScript(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateOracleScript) (*sdk.Result, error) {
	if k.IsNamePresent(ctx, types.OracleScriptStoreKeyString(msg.Name)) {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNameExists, "Name already exists")
	}
	k.AddOracleScriptFile(msg.Code, msg.Name)
	oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(msg.Name)
	// Get Data source and test case names from an oracle script
	aiDataSources, testCases, err := k.GetDNamesTcNames(oscriptPath)
	if err != nil {
		return nil, err
	}

	// collect minimum fees required to run the oracle script
	minimumFees, err := k.GetMinimumFees(ctx, aiDataSources, testCases)
	if err != nil {
		// erase because the script file is not properly added into the chain yet
		k.EraseOracleScriptFile(msg.Name)
		return nil, err
	}
	k.SetOracleScript(ctx, msg.Name, types.NewOracleScript(msg.Name, msg.Owner, msg.Description, minimumFees))
	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetOracleScript,
			sdk.NewAttribute(types.AttributeOracleScriptName, msg.Name),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgEditOracleScript is a function message for editing a oScript
func handleMsgEditOracleScript(ctx sdk.Context, k keeper.Keeper, msg types.MsgEditOracleScript) (*sdk.Result, error) {

	// Validate the oScript inputs
	oScript, err := k.GetOracleScript(ctx, msg.OldName)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}
	if !oScript.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrap(types.ErrEditorNotAuthorized, "Only owner can edit the oScript")
	}

	// Get Data source and test case names from an oracle script
	oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(msg.NewName)
	aiDataSources, testCases, err := k.GetDNamesTcNames(oscriptPath)
	if err != nil {
		return nil, err
	}

	// collect minimum fees required to run the oracle script
	minimumFees, err := k.GetMinimumFees(ctx, aiDataSources, testCases)
	if err != nil {
		return nil, err
	}

	oScript = types.NewOracleScript(msg.NewName, msg.Owner, msg.Description, minimumFees)

	k.EditOracleScript(ctx, msg.OldName, msg.NewName, oScript)
	if msg.OldName != msg.NewName {
		k.AddOracleScriptFile(msg.Code, msg.NewName)
	} else {
		k.EditOracleScriptFile(msg.Code, msg.NewName)
	}

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditOracleScript,
			sdk.NewAttribute(types.AttributeOracleScriptName, msg.NewName),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgCreateAIDataSource is a function message setting data source
func handleMsgCreateAIDataSource(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateAIDataSource) (*sdk.Result, error) {
	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.Fees)
	aiDataSource := types.NewAIDataSource(msg.Name, msg.Owner, fees, msg.Description)

	if k.IsNamePresent(ctx, types.DataSourceStoreKeyString(msg.Name)) {
		return nil, sdkerrors.Wrap(types.ErrDataSourceNameExists, "Name already exists")
	}

	k.SetAIDataSource(ctx, msg.Name, aiDataSource)
	k.AddAIDataSourceFile(msg.Code, msg.Name)

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetDataSource,
			sdk.NewAttribute(types.AttributeDataSourceName, msg.Name),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgEditAIDataSource is a function message for editing a data source
func handleMsgEditAIDataSource(ctx sdk.Context, k keeper.Keeper, msg types.MsgEditAIDataSource) (*sdk.Result, error) {

	aiDataSource, err := k.GetAIDataSource(ctx, msg.OldName)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}
	if !aiDataSource.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrap(types.ErrEditorNotAuthorized, "Only owner can edit the data source")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.Fees)
	aiDataSource = types.NewAIDataSource(msg.NewName, msg.Owner, fees, msg.Description)

	k.EditAIDataSource(ctx, msg.OldName, msg.NewName, aiDataSource)
	if msg.OldName != msg.NewName {
		k.AddAIDataSourceFile(msg.Code, msg.NewName)
	} else {
		k.EditAIDataSourceFile(msg.Code, msg.NewName)
	}

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditDataSource,
			sdk.NewAttribute(types.AttributeDataSourceName, msg.NewName),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
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
	// We need to calculate the final 60% fee given by the user because the remaining 40% must be reserved for the proposer and validators.
	providedCoins, _ := sdk.NewDecCoinsFromCoins(fees...).MulDecTruncate(rewardRatio).TruncateDecimal()

	//validatorFees, err := k.GetValidatorFees(ctx, providedCoins, validators)

	// collect data source name from the oScript script
	oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(msg.MsgAIRequest.OracleScriptName)

	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.GetDNamesTcNames(oscriptPath)
	if err != nil {
		return nil, err
	}

	//var finalResult int

	var testcaseObjs []TestCase

	var dataSourceObjs []AIDataSource

	//var testCaseResults []types.TestCaseResult

	var totalFees sdk.Coins

	// we have different test cases, so we need to loop through them
	for i := 0; i < len(testCases); i++ {
		// loop to run the test case
		// collect all the test cases object to store in the ai request
		testCase, err := k.GetTestCase(ctx, testCases[i])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, fmt.Sprintf("failed to get test case: %s", err.Error()))
		}
		testcaseObjs = append(testcaseObjs, testCase)
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(testCase.Fees...)
	}

	for j := 0; j < len(aiDataSources); j++ {
		// collect all the data source objects to store in the ai request
		aiDataSource, err := k.GetAIDataSource(ctx, aiDataSources[j])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, fmt.Sprintf("failed to get data source: %s", err.Error()))
		}
		dataSourceObjs = append(dataSourceObjs, aiDataSource)
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(aiDataSource.Fees...)
	}

	// If the total fee is larger than the fee provided by the user then we return error
	if totalFees.IsAnyGT(providedCoins) {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Fees given by the users are less than the total fees needed")
	}

	// set a new request with the aggregated result into blockchain
	request := types.NewAIRequest(msg.MsgAIRequest.RequestID, msg.MsgAIRequest.OracleScriptName, msg.MsgAIRequest.Creator, validators, ctx.BlockHeight(), dataSourceObjs, testcaseObjs, fees, msg.MsgAIRequest.Input, msg.MsgAIRequest.ExpectedOutput)

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

// handleMsgSetKYCRequest is a function message setting a new AI request
func handleMsgSetKYCRequest(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetKYCRequest) (*sdk.Result, error) {

	validators, err := k.RandomValidators(ctx, msg.MsgAIRequest.ValidatorCount, []byte(msg.MsgAIRequest.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.MsgAIRequest.Fees)
	// Compute the fee allocated for oracle module to distribute to active validators.
	rewardRatio := sdk.NewDecWithPrec(int64(k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)), 2)
	// We need to calculate the final 60% fee given by the user because the remaining 40% must be reserved for the proposer and validators.
	providedCoins, _ := sdk.NewDecCoinsFromCoins(fees...).MulDecTruncate(rewardRatio).TruncateDecimal()

	//validatorFees, err := k.GetValidatorFees(ctx, providedCoins, validators)

	// collect data source name from the oScript script
	oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(msg.MsgAIRequest.OracleScriptName)

	//use "data source" as an argument to collect the data source script name
	cmd := exec.Command("bash", oscriptPath, "aiDataSource")
	cmd.Stdin = strings.NewReader("some input")
	var dataSourceName bytes.Buffer
	cmd.Stdout = &dataSourceName
	err = cmd.Run()
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrFailedToOpenFile, "failed to collect data source name")
	}

	// collect data source result from the script
	result := strings.TrimSuffix(dataSourceName.String(), "\n")

	aiDataSources := strings.Fields(result)

	//use "test case" as an argument to collect the test case script name
	cmd = exec.Command("bash", oscriptPath, "testcase")
	cmd.Stdin = strings.NewReader("some input")
	var testCaseName bytes.Buffer
	cmd.Stdout = &testCaseName
	err = cmd.Run()
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrFailedToOpenFile, fmt.Sprintf("failed to collect test case name: %s", result))
	}

	// collect data source result from the script
	result = strings.TrimSuffix(testCaseName.String(), "\n")

	testCases := strings.Fields(result)

	//var finalResult int

	var testcaseObjs []TestCase

	var dataSourceObjs []AIDataSource

	//var testCaseResults []types.TestCaseResult

	var totalFees sdk.Coins

	// we have different test cases, so we need to loop through them
	for i := 0; i < len(testCases); i++ {
		// loop to run the test case
		// collect all the test cases object to store in the ai request
		testCase, err := k.GetTestCase(ctx, testCases[i])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, fmt.Sprintf("failed to get test case: %s", err.Error()))
		}
		testcaseObjs = append(testcaseObjs, testCase)
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(testCase.Fees...)
	}

	for j := 0; j < len(aiDataSources); j++ {
		// collect all the data source objects to store in the ai request
		aiDataSource, err := k.GetAIDataSource(ctx, aiDataSources[j])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, fmt.Sprintf("failed to get data source: %s", err.Error()))
		}
		dataSourceObjs = append(dataSourceObjs, aiDataSource)
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(aiDataSource.Fees...)
	}

	// If the total fee is larger than the fee provided by the user then we return error
	if totalFees.IsAnyGT(providedCoins) {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Fees given by the users are less than the total fees needed")
	}

	// set a new request with the aggregated result into blockchain
	request := types.NewAIRequest(msg.MsgAIRequest.RequestID, msg.MsgAIRequest.OracleScriptName, msg.MsgAIRequest.Creator, validators, ctx.BlockHeight(), dataSourceObjs, testcaseObjs, fees, msg.MsgAIRequest.Input, msg.MsgAIRequest.ExpectedOutput)

	//fmt.Printf("request result: %s\n", outOracleScript.String())

	k.SetAIRequest(ctx, request.RequestID, request)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeSetKYCRequest)
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
		sdk.NewAttribute(types.AttributeRequestImageHash, msg.ImageHash),
		sdk.NewAttribute(types.AttributeRequestImageName, msg.ImageName),
		sdk.NewAttribute(types.AttributeRequestValidatorCount, string(msg.MsgAIRequest.ValidatorCount)),
		sdk.NewAttribute(types.AttributeRequestInput, msg.MsgAIRequest.Input),
		sdk.NewAttribute(types.AttributeRequestExpectedOutput, msg.MsgAIRequest.ExpectedOutput),
	)
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgCreateTestCase is a function message setting test case
func handleMsgCreateTestCase(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateTestCase) (*sdk.Result, error) {
	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.Fees)
	testCase := types.NewTestCase(msg.Name, msg.Owner, fees, msg.Description)

	if k.IsNamePresent(ctx, types.TestCaseStoreKeyString(msg.Name)) {
		return nil, sdkerrors.Wrap(types.ErrTestCaseNameExists, "Name already exists")
	}

	k.CreateTestCase(ctx, msg.Name, testCase)
	k.AddTestCaseFile(msg.Code, msg.Name)

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTestCase,
			sdk.NewAttribute(types.AttributeTestCaseName, msg.Name),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgEditTestCase is a function message for editing a test case
func handleMsgEditTestCase(ctx sdk.Context, k keeper.Keeper, msg types.MsgEditTestCase) (*sdk.Result, error) {

	testCase, err := k.GetTestCase(ctx, msg.OldName)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, err.Error())
	}
	if !testCase.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrap(types.ErrEditorNotAuthorized, "Only owner can edit the data source")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.Fees)
	testCase = types.NewTestCase(msg.NewName, msg.Owner, fees, msg.Description)

	k.EditTestCase(ctx, msg.OldName, msg.NewName, testCase)
	if msg.OldName != msg.NewName {
		k.AddTestCaseFile(msg.Code, msg.NewName)
	} else {
		k.EditTestCaseFile(msg.Code, msg.NewName)
	}

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditDataSource,
			sdk.NewAttribute(types.AttributeDataSourceName, msg.NewName),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// this handler will be triggered when the websocket create a MsgCreateReport
func handleMsgAddReport(ctx sdk.Context, k Keeper, msg types.MsgCreateReport) (*sdk.Result, error) {
	validator := types.NewValidator(msg.Validator, k.GetValidator(ctx, msg.Validator).GetConsensusPower(), "active")
	report := types.NewReport(msg.RequestID, validator, msg.DataSourceResults, msg.TestCaseResults, ctx.BlockHeight(), msg.Fees, msg.AggregatedResult)
	err := k.AddReport(ctx, msg.RequestID, report)
	if err != nil {
		return nil, err
	}
	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTestCase,
			sdk.NewAttribute(types.AttributeTestCaseName, msg.Reporter.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgAddReporter(ctx sdk.Context, k Keeper, m types.MsgAddReporter) (*sdk.Result, error) {
	err := k.AddReporter(ctx, m.Validator, m.Reporter)
	if err != nil {
		return nil, err
	}
	// // TODO: Define your msg events
	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		types.EventTypeCreateTestCase,
	// 		sdk.NewAttribute(types.AttributeTestCaseName, m.Reporter.String()),
	// 	),
	// )
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveReporter(ctx sdk.Context, k Keeper, m types.MsgRemoveReporter) (*sdk.Result, error) {
	err := k.RemoveReporter(ctx, m.Validator, m.Reporter)
	if err != nil {
		return nil, err
	}
	// ctx.EventManager().EmitEvent(sdk.NewEvent(
	// 	types.EventTypeRemoveReporter,
	// 	sdk.NewAttribute(types.AttributeKeyValidator, m.Validator.String()),
	// 	sdk.NewAttribute(types.AttributeKeyReporter, m.Reporter.String()),
	// ))
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateStrategy(ctx sdk.Context, k Keeper, msg types.MsgCreateStrategy) (*sdk.Result, error) {

	strategy := types.NewStrategy(msg.StratID, msg.StratName, msg.StratFlow, msg.PerformanceFee, msg.PerformanceMax, msg.WithdrawalFee, msg.WithdrawalMax, msg.GovernanceAddr, msg.StrategistAddr)

	k.CreateStrategy(ctx, strategy.StratName, strategy)
	// // TODO: Define your msg events
	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		types.EventTypeCreateTestCase,
	// 		sdk.NewAttribute(types.AttributeTestCaseName, m.Reporter.String()),
	// 	),
	// )
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
