package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateAIDataSource(goCtx context.Context, msg *types.MsgCreateAIDataSource) (*types.MsgCreateAIDataSourceRes, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoinsNormalized(msg.Fees)
	aiDataSource := types.NewAIDataSource(msg.Name, msg.Owner, fees, msg.Description)

	fmt.Println("ai data source: ", aiDataSource)

	//k.SetAIDataSource(ctx, msg.Name, aiDataSource)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "abcd"),
		),
	)

	return &types.MsgCreateAIDataSourceRes{}, nil
}

// Edit an existing data source
func (k msgServer) EditAIDataSource(ctx context.Context, in *types.MsgEditAIDataSource) (*types.MsgEditAIDataSourceRes, error) {
	return &types.MsgEditAIDataSourceRes{}, nil
}

// Create a new oracle script

// Edit an existing oracle script
func (k msgServer) EditOracleScript(ctx context.Context, in *types.MsgEditOracleScript) (*types.MsgEditOracleScriptRes, error) {
	return &types.MsgEditOracleScriptRes{}, nil
}

// Create a new test case
func (k msgServer) CreateTestCase(ctx context.Context, in *types.MsgCreateTestCase) (*types.MsgCreateTestCaseRes, error) {
	return &types.MsgCreateTestCaseRes{}, nil
}

// Edit an existing test case
func (k msgServer) EditTestCase(ctx context.Context, in *types.MsgEditTestCase) (*types.MsgEditTestCaseRes, error) {
	return &types.MsgEditTestCaseRes{}, nil
}

func (k msgServer) CreateOracleScript(ctx context.Context, in *types.MsgCreateOracleScript) (*types.MsgCreateOracleScriptRes, error) {
	return &types.MsgCreateOracleScriptRes{}, nil
}
