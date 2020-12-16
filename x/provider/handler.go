package provider

import (
	"fmt"

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
		case types.MsgCreateTestCase:
			return handleMsgCreateTestCase(ctx, k, msg)
		case types.MsgEditTestCase:
			return handleMsgEditTestCase(ctx, k, msg)
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

	// collect minimum fees required to run the oracle script (for 1 validator)
	minimumFees, err := k.GetMinimumFees(ctx, msg.DataSources, msg.TestCases, 1)
	if err != nil {
		return nil, err
	}
	k.SetOracleScript(ctx, msg.Name, types.NewOracleScript(msg.Name, msg.Owner, msg.Description, minimumFees, msg.DataSources, msg.TestCases))
	k.AddOracleScriptFile(msg.Code, msg.Name)
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

	// collect minimum fees required to run the oracle script (for one validator)
	minimumFees, err := k.GetMinimumFees(ctx, msg.DataSources, msg.TestCases, 1)
	if err != nil {
		return nil, err
	}

	oScript = types.NewOracleScript(msg.NewName, msg.Owner, msg.Description, minimumFees, msg.DataSources, msg.TestCases)

	k.EditOracleScript(ctx, msg.OldName, msg.NewName, msg.Code, oScript)

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

	k.EditAIDataSource(ctx, msg.OldName, msg.NewName, msg.Code, aiDataSource)

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditDataSource,
			sdk.NewAttribute(types.AttributeDataSourceName, msg.NewName),
		),
	)

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

	k.SetTestCase(ctx, msg.Name, testCase)
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

	k.EditTestCase(ctx, msg.OldName, msg.NewName, msg.Code, testCase)

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditTestCase,
			sdk.NewAttribute(types.AttributeTestCaseName, msg.NewName),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
