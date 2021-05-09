package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/aioracle/types"
	aioracle "github.com/oraichain/orai/x/aioracle/types"
)

type msgServer struct {
	querier *Querier
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(querier *Querier) types.MsgServer {
	return &msgServer{querier: querier}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateAIOracle(goCtx context.Context, msg *types.MsgSetAIOracleReq) (*types.MsgSetAIOracleRes, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators, err := k.querier.keeper.RandomValidators(ctx, int(msg.ValidatorCount), []byte(msg.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}
	// we can safely parse to acc address because we have validated them
	contract, _ := sdk.AccAddressFromBech32(msg.Contract)
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	// validate if the request id exists or not
	if k.querier.keeper.HasAIOracle(ctx, msg.RequestID) {
		return nil, sdkerrors.Wrap(types.ErrRequestInvalid, "The request id already exists")
	}

	// check size of the request
	maxBytes := int(k.querier.keeper.GetParam(ctx, types.KeyMaximumAIOracleReqBytes))
	// threshold for the size of the request
	if len(msg.Input) > maxBytes {
		return nil, sdkerrors.Wrap(types.ErrRequestInvalid, "The request is too large")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	providedFees, _ := sdk.ParseCoinsNormalized(msg.Fees)

	requiredFees, err := k.querier.keeper.calculateMinimumFees(ctx, msg.TestOnly, contract, len(validators))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrQueryMinFees, fmt.Sprintf("Error getting minimum fees from oracle script with err: %v", err))
	}

	// If the total fee is larger than the fee provided by the user then we return error
	if requiredFees.IsAnyGT(providedFees) {
		k.querier.keeper.Logger(ctx).Error(fmt.Sprintf("Your payment fees is less than required\n"))
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, fmt.Sprintf("Fees given: %v, where fees required is: %v", providedFees, requiredFees))
	}

	// check if the account has enough spendable coins
	spendableCoins := k.querier.keeper.BankKeeper.SpendableCoins(ctx, creator)
	// If the total fee is larger or equal to the spendable coins of the user then we return error
	if providedFees.IsAnyGTE(spendableCoins) {
		k.querier.keeper.Logger(ctx).Error(fmt.Sprintf("Your account has run out of tokens to create the AI Request\n"))
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Your account has run out of tokens to create the AI Request")
	}

	// substract coins in the creator wallet to charge fees
	err = k.querier.keeper.BankKeeper.SubtractCoins(ctx, creator, providedFees)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, fmt.Sprintf("Your account has run out of tokens to create the AI Request, or there is something wrong with error: %v", err))
	}

	// set a new request with the aggregated result into blockchain
	request := types.NewAIOracle(msg.RequestID, contract, creator, validators, ctx.BlockHeight(), providedFees, msg.Input, msg.TestOnly)

	k.querier.keeper.SetAIOracle(ctx, request.RequestID, request)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrRequestInvalid, fmt.Sprintf("Cannot marshal request %v with error: %v", request, err))
	}
	// Emit an event describing a data request and asked validators.
	eventKey := types.AttributeReport
	if msg.TestOnly == true {
		eventKey = types.AttributeBaseReport
	}
	event := sdk.NewEvent(types.EventTypeRequestWithData)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequest, string(requestBytes)),
		sdk.NewAttribute(eventKey, "report types"),
	)

	ctx.EventManager().EmitEvent(event)

	return types.NewMsgSetAIOracleRes(
		request.GetRequestID(), msg.GetContract(),
		msg.GetCreator(), request.GetFees().String(), msg.GetValidatorCount(),
		request.GetInput(), request.GetTestOnly(),
	), nil
}

func (k *Keeper) calculateMinimumFees(ctx sdk.Context, isTestOnly bool, contractAddr sdk.AccAddress, numVal int) (sdk.Coins, error) {
	querier := NewQuerier(k)
	goCtx := sdk.WrapSDKContext(ctx)
	entries := &types.ResponseEntryPointContract{}
	var err error
	if isTestOnly {
		entries, err = querier.TestCaseEntries(goCtx, &types.QueryTestCaseEntriesContract{
			Contract: contractAddr,
			Request:  &types.EmptyParams{},
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err = querier.DataSourceEntries(goCtx, &types.QueryDataSourceEntriesContract{
			Contract: contractAddr,
			Request:  &types.EmptyParams{},
		})
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("entries: ", entries)
	minFees := sdk.Coins{}
	for _, entry := range entries.GetData() {
		minFees = minFees.Add(entry.GetProviderFees()...)
	}
	valFees := k.CalculateValidatorFees(ctx, minFees)
	minFees = minFees.Add(valFees...)
	// since there are more than 1 validator, we need to multiply those fees
	minFees, _ = sdk.NewDecCoinsFromCoins(minFees...).MulDec(sdk.NewDec(int64(numVal))).TruncateDecimal()
	return minFees, nil
}

// this handler will be triggered when the websocket create a MsgCreateReport
func (m msgServer) CreateReport(goCtx context.Context, msg *types.MsgCreateReport) (*types.MsgCreateReportRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if m.querier.keeper.HasReport(ctx, msg.BaseReport.GetRequestId(), msg.BaseReport.GetValidatorAddress()) {
		return nil, fmt.Errorf("Error: Validator already reported")
	}

	request, err := m.querier.keeper.GetAIOracle(ctx, msg.BaseReport.GetRequestId())
	if err != nil {
		return nil, fmt.Errorf("Error: Cannot find AI request")
	}

	if !m.querier.ContainsVal(request.GetValidators(), msg.BaseReport.GetValidatorAddress()) {
		return nil, sdkerrors.Wrap(types.ErrValidatorNotFound, fmt.Sprintln("Failed to find the requested validator"))
	}

	report := types.NewReport(msg.BaseReport.GetRequestId(), msg.DataSourceResults, ctx.BlockHeight(), msg.AggregatedResult, msg.BaseReport.GetValidatorAddress(), msg.ResultStatus, msg.BaseReport.GetFees())

	if !m.querier.validateReportBasic(ctx, request, report) {
		return nil, sdkerrors.Wrap(types.ErrMsgReportInvalid, "Report does not pass the validation step")
	}

	// call keeper
	err = m.querier.keeper.SetReport(ctx, msg.BaseReport.GetRequestId(), report)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateReportRes{
		BaseReport: &aioracle.BaseReport{
			RequestId:        msg.BaseReport.GetRequestId(),
			ValidatorAddress: msg.BaseReport.ValidatorAddress,
			Fees:             msg.BaseReport.Fees,
		},
		DataSourceResults: msg.GetDataSourceResults(),
		AggregatedResult:  msg.GetAggregatedResult(),
		ResultStatus:      msg.GetResultStatus(),
	}, nil
}

// this handler will be triggered when the websocket create a MsgCreateReport
func (m msgServer) CreateTestCaseReport(goCtx context.Context, msg *types.MsgCreateTestCaseReport) (*types.MsgCreateTestCaseReportRes, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	if m.querier.keeper.HasTestCaseReport(ctx, msg.BaseReport.GetRequestId(), msg.BaseReport.GetValidatorAddress()) {
		return nil, fmt.Errorf("Error: Validator already reported")
	}

	request, err := m.querier.keeper.GetAIOracle(ctx, msg.BaseReport.GetRequestId())
	if err != nil {
		return nil, fmt.Errorf("Error: Cannot find AI request")
	}

	if !m.querier.ContainsVal(request.GetValidators(), msg.BaseReport.GetValidatorAddress()) {
		return nil, sdkerrors.Wrap(types.ErrValidatorNotFound, fmt.Sprintln("Failed to find the requested validator"))
	}

	report := types.NewTestCaseReport(msg.BaseReport.GetRequestId(), msg.GetResultsWithTestCase(), ctx.BlockHeight(), msg.BaseReport.GetValidatorAddress(), msg.BaseReport.GetFees())

	if !m.querier.validateTestCaseReportBasic(ctx, request, report) {
		return nil, sdkerrors.Wrap(types.ErrMsgReportInvalid, "Test case report does not pass the validation step")
	}

	// call keeper
	err = m.querier.keeper.SetTestCaseReport(ctx, msg.BaseReport.GetRequestId(), report)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateTestCaseReportRes{
		BaseReport: &aioracle.BaseReport{
			RequestId:        msg.BaseReport.GetRequestId(),
			ValidatorAddress: msg.BaseReport.ValidatorAddress,
			Fees:             msg.BaseReport.Fees,
		},
		ResultsWithTestCase: msg.GetResultsWithTestCase(),
	}, nil
}
