package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/types"
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

// this handler will be triggered when the websocket create a MsgCreateReport
func (m msgServer) AddReport(goCtx context.Context, msg *types.MsgCreateReport) (*types.MsgCreateReportRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// basic validation before adding the report
	if m.keeper.HasReport(ctx, msg.RequestID, msg.Reporter.Validator) {
		return nil, fmt.Errorf("Error: Validator already reported")
	}

	request, err := m.keeper.aiRequestKeeper.GetAIRequest(ctx, msg.RequestID)
	if err != nil {
		return nil, fmt.Errorf("Error: Cannot find AI request")
	}

	if m.keeper.ValidateReport(ctx, msg.GetReporter(), request) != nil {
		return nil, fmt.Errorf("Error: cannot find reporter in the AI request")
	}

	report := types.NewReport(msg.RequestID, msg.DataSourceResults, msg.TestCaseResults, ctx.BlockHeight(), msg.Fees, msg.AggregatedResult, msg.Reporter, msg.ResultStatus)

	// call keeper
	err = m.keeper.AddReport(ctx, msg.RequestID, report)
	if err != nil {
		return nil, err
	}
	// Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetReport,
			sdk.NewAttribute(types.AttributeReport, msg.Reporter.String()),
		),
	)

	return &types.MsgCreateReportRes{
		RequestID:         msg.GetRequestID(),
		DataSourceResults: msg.GetDataSourceResults(),
		TestCaseResults:   msg.GetTestCaseResults(),
		Reporter:          msg.GetReporter(),
		Fees:              msg.GetFees(),
		AggregatedResult:  msg.GetAggregatedResult(),
		ResultStatus:      msg.GetResultStatus(),
	}, nil
}
