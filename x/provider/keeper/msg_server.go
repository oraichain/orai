package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/provider/types"
)

type msgServer struct {
	keeper *Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CreateAIDataSource(goCtx context.Context, msg *types.MsgCreateAIDataSource) (*types.MsgCreateAIDataSourceRes, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if m.keeper.HasDataSource(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrDataSourceNameExists, "Name already exists")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoinsNormalized(msg.Fees)
	aiDataSource := types.NewAIDataSource(msg.Name, msg.Contract, msg.Owner, fees, msg.Description)

	err := m.keeper.SetAIDataSource(ctx, msg.Name, aiDataSource)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotSetDataSource, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetDataSource,
			sdk.NewAttribute(types.AttributeDataSourceName, msg.Name),
			sdk.NewAttribute(types.AttributeContractAddress, msg.Contract),
		),
	)

	return &types.MsgCreateAIDataSourceRes{
		Name:        msg.GetName(),
		Description: msg.GetDescription(),
		Contract:    msg.GetContract(),
		Owner:       msg.GetOwner(),
		Fees:        msg.GetFees(),
	}, nil
}

// Edit an existing data source
func (m msgServer) EditAIDataSource(goCtx context.Context, msg *types.MsgEditAIDataSource) (*types.MsgEditAIDataSourceRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.keeper.HasDataSource(ctx, msg.OldName) {
		return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, "Data source is not found")
	}

	aiDataSource, err := m.keeper.GetAIDataSource(ctx, msg.OldName)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}
	if !aiDataSource.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrap(types.ErrEditorNotAuthorized, "Only owner can edit the data source")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoinsNormalized(msg.Fees)

	aiDataSource = types.NewAIDataSource(msg.NewName, msg.Contract, msg.Owner, fees, msg.Description)

	err = m.keeper.EditAIDataSource(ctx, msg.OldName, msg.NewName, aiDataSource)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotSetDataSource, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditDataSource,
			sdk.NewAttribute(types.AttributeDataSourceName, msg.NewName),
			sdk.NewAttribute(types.AttributeContractAddress, msg.Contract),
		),
	)

	return &types.MsgEditAIDataSourceRes{
		Name:        msg.GetNewName(),
		Description: msg.GetDescription(),
		Contract:    msg.GetContract(),
		Owner:       msg.GetOwner(),
		Fees:        msg.GetFees(),
	}, nil
}

// CreateOracleScript: Create a new oracle script
func (m msgServer) CreateOracleScript(goCtx context.Context, msg *types.MsgCreateOracleScript) (*types.MsgCreateOracleScriptRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if m.keeper.HasOracleScript(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNameExists, "Name already exists")
	}

	// collect oracle script reward percentage
	rewardPercentage := m.keeper.GetOracleScriptRewardPercentageParam(ctx)
	// collect minimum fees required to run the oracle script (for 1 validator)
	minimumFees, err := m.keeper.GetMinimumFees(ctx, msg.DataSources, msg.TestCases, 1, rewardPercentage)
	if err != nil {
		return nil, err
	}

	err = m.keeper.SetOracleScript(ctx, msg.Name,
		types.NewOracleScript(msg.Name, msg.Contract, msg.Owner,
			msg.Description, minimumFees, msg.DataSources, msg.TestCases))

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotSetOracleScript, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetOracleScript,
			sdk.NewAttribute(types.AttributeOracleScriptName, msg.Name),
			sdk.NewAttribute(types.AttributeContractAddress, msg.Contract),
		),
	)

	return &types.MsgCreateOracleScriptRes{
		Name:        msg.GetName(),
		Description: msg.GetDescription(),
		Contract:    msg.GetContract(),
		Owner:       msg.GetOwner(),
		Fees:        msg.GetFees(),
		DataSources: msg.GetDataSources(),
		TestCases:   msg.GetTestCases(),
	}, nil
}

// EditOracleScript: Edit an existing oracle script
func (m msgServer) EditOracleScript(goCtx context.Context, msg *types.MsgEditOracleScript) (*types.MsgEditOracleScriptRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.keeper.HasOracleScript(ctx, msg.OldName) {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNotFound, "Oracle script does not exist")
	}

	// Validate the oScript inputs
	oScript, err := m.keeper.GetOracleScript(ctx, msg.OldName)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}
	if !oScript.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrap(types.ErrEditorNotAuthorized, "Only owner can edit the oScript")
	}

	// collect oracle script reward percentage
	rewardPercentage := m.keeper.GetOracleScriptRewardPercentageParam(ctx)
	// collect minimum fees required to run the oracle script (for one validator)
	minimumFees, err := m.keeper.GetMinimumFees(ctx, msg.DataSources, msg.TestCases, 1, rewardPercentage)
	if err != nil {
		return nil, err
	}

	oScript = types.NewOracleScript(msg.NewName, msg.Contract, msg.Owner,
		msg.Description, minimumFees, msg.DataSources, msg.TestCases)

	err = m.keeper.EditOracleScript(ctx, msg.OldName, msg.NewName, oScript)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotSetOracleScript, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditOracleScript,
			sdk.NewAttribute(types.AttributeOracleScriptName, msg.NewName),
			sdk.NewAttribute(types.AttributeContractAddress, msg.Contract),
		),
	)

	return &types.MsgEditOracleScriptRes{
		Name:        msg.GetNewName(),
		Description: msg.GetDescription(),
		Contract:    msg.GetContract(),
		Owner:       msg.GetOwner(),
		Fees:        msg.GetFees(),
		DataSources: msg.GetDataSources(),
		TestCases:   msg.GetTestCases(),
	}, nil
}

// CreateTestCase: Create a new test case
func (m msgServer) CreateTestCase(goCtx context.Context, msg *types.MsgCreateTestCase) (*types.MsgCreateTestCaseRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if m.keeper.HasTestCase(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrTestCaseNameExists, "Name already exists")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoinsNormalized(msg.Fees)

	testCase := types.NewTestCase(msg.Name, msg.Contract, msg.Owner, fees, msg.Description)

	err := m.keeper.SetTestCase(ctx, msg.Name, testCase)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotSetTestCase, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTestCase,
			sdk.NewAttribute(types.AttributeTestCaseName, msg.Name),
			sdk.NewAttribute(types.AttributeContractAddress, msg.Contract),
		),
	)

	return &types.MsgCreateTestCaseRes{
		Name:        msg.GetName(),
		Description: msg.GetDescription(),
		Contract:    msg.GetContract(),
		Owner:       msg.GetOwner(),
		Fees:        msg.GetFees(),
	}, nil
}

// EditTestCase: Edit an existing test case
func (m msgServer) EditTestCase(goCtx context.Context, msg *types.MsgEditTestCase) (*types.MsgEditTestCaseRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.keeper.HasTestCase(ctx, msg.NewName) {
		return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, "Test case does not exist")
	}

	testCase, err := m.keeper.GetTestCase(ctx, msg.OldName)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, err.Error())
	}
	if !testCase.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrap(types.ErrEditorNotAuthorized, "Only owner can edit the data source")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoinsNormalized(msg.Fees)
	testCase = types.NewTestCase(msg.NewName, msg.Contract, msg.Owner, fees, msg.Description)

	err = m.keeper.EditTestCase(ctx, msg.OldName, msg.NewName, testCase)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotSetTestCase, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEditTestCase,
			sdk.NewAttribute(types.AttributeTestCaseName, msg.NewName),
			sdk.NewAttribute(types.AttributeContractAddress, msg.Contract),
		),
	)

	return &types.MsgEditTestCaseRes{
		Name:        msg.GetNewName(),
		Description: msg.GetDescription(),
		Contract:    msg.GetContract(),
		Owner:       msg.GetOwner(),
		Fees:        msg.GetFees(),
	}, nil
}
