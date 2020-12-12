package websocket

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/websocket/keeper"
	"github.com/oraichain/orai/x/websocket/types"
)

// NewHandler creates an sdk.Handler for all the provider type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		// case types.MsgSetKYCRequest:
		// 	return handleMsgSetKYCRequest(ctx, k, msg)
		case types.MsgCreateReport:
			return handleMsgAddReport(ctx, k, msg)
		case types.MsgAddReporter:
			return handleMsgAddReporter(ctx, k, msg)
		case types.MsgRemoveReporter:
			return handleMsgRemoveReporter(ctx, k, msg)
		// case types.MsgCreateStrategy:
		// 	return handleMsgCreateStrategy(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// this handler will be triggered when the websocket create a MsgCreateReport
func handleMsgAddReport(ctx sdk.Context, k Keeper, msg types.MsgCreateReport) (*sdk.Result, error) {
	// validator := types.NewValidator(msg.Reporter.Validator, k.GetValidator(ctx, msg.Reporter.Validator).GetConsensusPower(), "active")
	report := types.NewReport(msg.RequestID, msg.DataSourceResults, msg.TestCaseResults, ctx.BlockHeight(), msg.Fees, msg.AggregatedResult, msg.Reporter, msg.ResultStatus)
	// basic validation before adding the report
	if k.HasReport(ctx, msg.RequestID, msg.Reporter.Validator) {
		return nil, fmt.Errorf("Error: Validator already reported")
	}
	err := k.AddReport(ctx, msg.RequestID, report)
	if err != nil {
		return nil, err
	}
	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetReport,
			sdk.NewAttribute(types.AttributeReport, msg.Reporter.String()),
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

// func handleMsgCreateStrategy(ctx sdk.Context, k Keeper, msg types.MsgCreateStrategy) (*sdk.Result, error) {

// 	strategy := types.NewStrategy(msg.StratID, msg.StratName, msg.StratFlow, msg.PerformanceFee, msg.PerformanceMax, msg.WithdrawalFee, msg.WithdrawalMax, msg.GovernanceAddr, msg.StrategistAddr)

// 	k.CreateStrategy(ctx, strategy.StratName, strategy)
// 	// // TODO: Define your msg events
// 	// ctx.EventManager().EmitEvent(
// 	// 	sdk.NewEvent(
// 	// 		types.EventTypeCreateTestCase,
// 	// 		sdk.NewAttribute(types.AttributeTestCaseName, m.Reporter.String()),
// 	// 	),
// 	// )
// 	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
// }
