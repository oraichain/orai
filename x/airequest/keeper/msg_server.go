package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/types"
	airequest "github.com/oraichain/orai/x/airequest/types"
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

func (k msgServer) CreateAIRequest(goCtx context.Context, msg *types.MsgSetAIRequestReq) (*types.MsgSetAIRequestRes, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators, err := k.keeper.RandomValidators(ctx, int(msg.ValidatorCount), []byte(msg.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}
	// we can safely parse to acc address because we have validated them
	contract, _ := sdk.AccAddressFromBech32(msg.Contract)
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	// validate if the request id exists or not
	if k.keeper.HasAIRequest(ctx, msg.RequestID) {
		return nil, sdkerrors.Wrap(types.ErrRequestInvalid, "The request id already exists")
	}

	// check size of the request
	maxBytes := int(k.keeper.GetParam(ctx, types.KeyMaximumAIRequestReqBytes))
	// threshold for the size of the request
	if len(msg.Input) > maxBytes {
		return nil, sdkerrors.Wrap(types.ErrRequestInvalid, "The request is too large")
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	providedFees, _ := sdk.ParseCoinsNormalized(msg.Fees)

	requiredFees, err := k.keeper.calculateMinimumFees(ctx, msg.TestOnly, contract, len(validators))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrQueryMinFees, fmt.Sprintf("Error getting minimum fees from oracle script with err: %v", err))
	}

	// If the total fee is larger than the fee provided by the user then we return error
	if requiredFees.IsAnyGT(providedFees) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Your payment fees is less than required\n"))
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, fmt.Sprintf("Fees given: %v, where fees required is: %v", providedFees, requiredFees))
	}

	// check if the account has enough spendable coins
	spendableCoins := k.keeper.BankKeeper.SpendableCoins(ctx, creator)
	// If the total fee is larger or equal to the spendable coins of the user then we return error
	if providedFees.IsAnyGTE(spendableCoins) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Your account has run out of tokens to create the AI Request\n"))
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Your account has run out of tokens to create the AI Request")
	}

	// substract coins in the creator wallet to charge fees
	err = k.keeper.BankKeeper.SubtractCoins(ctx, creator, providedFees)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, fmt.Sprintf("Your account has run out of tokens to create the AI Request, or there is something wrong with error: %v", err))
	}

	// set a new request with the aggregated result into blockchain
	request := types.NewAIRequest(msg.RequestID, contract, creator, validators, ctx.BlockHeight(), providedFees, msg.Input, msg.TestOnly)

	k.keeper.SetAIRequest(ctx, request.RequestID, request)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequestWithData)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
		sdk.NewAttribute(types.AttributeContract, request.Contract.String()),
		sdk.NewAttribute(types.AttributeRequestCreator, creator.String()),
		sdk.NewAttribute(types.AttributeRequestValidatorCount, fmt.Sprint(msg.ValidatorCount)),
		sdk.NewAttribute(types.AttributeRequestInput, string(msg.Input)),
	)

	for _, validator := range validators {
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeRequestValidator, validator.String()),
		)
	}

	ctx.EventManager().EmitEvent(event)

	return types.NewMsgSetAIRequestRes(
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

	if m.keeper.HasReport(ctx, msg.BaseReport.GetRequestId(), msg.BaseReport.GetValidatorAddress()) {
		return nil, fmt.Errorf("Error: Validator already reported")
	}

	request, err := m.keeper.GetAIRequest(ctx, msg.BaseReport.GetRequestId())
	if err != nil {
		return nil, fmt.Errorf("Error: Cannot find AI request")
	}

	if m.keeper.ValidateReport(ctx, msg.BaseReport.GetValidatorAddress(), request) != nil {
		return nil, fmt.Errorf("Error: cannot find reporter in the AI request")
	}

	report := types.NewReport(msg.BaseReport.GetRequestId(), msg.DataSourceResults, ctx.BlockHeight(), msg.AggregatedResult, msg.BaseReport.GetValidatorAddress(), msg.ResultStatus, msg.BaseReport.GetFees())

	// call keeper
	err = m.keeper.SetReport(ctx, msg.BaseReport.GetRequestId(), report)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateReportRes{
		BaseReport: &airequest.BaseReport{
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
	if m.keeper.HasTestCaseReport(ctx, msg.BaseReport.GetRequestId(), msg.BaseReport.GetValidatorAddress()) {
		return nil, fmt.Errorf("Error: Validator already reported")
	}

	request, err := m.keeper.GetAIRequest(ctx, msg.BaseReport.GetRequestId())
	if err != nil {
		return nil, fmt.Errorf("Error: Cannot find AI request")
	}

	if m.keeper.ValidateReport(ctx, msg.BaseReport.GetValidatorAddress(), request) != nil {
		return nil, fmt.Errorf("Error: cannot find reporter in the AI request")
	}

	report := types.NewTestCaseReport(msg.BaseReport.GetRequestId(), msg.GetResultsWithTestCase(), ctx.BlockHeight(), msg.BaseReport.GetValidatorAddress(), msg.BaseReport.GetFees())

	// call keeper
	err = m.keeper.SetTestCaseReport(ctx, msg.BaseReport.GetRequestId(), report)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateTestCaseReportRes{
		BaseReport: &airequest.BaseReport{
			RequestId:        msg.BaseReport.GetRequestId(),
			ValidatorAddress: msg.BaseReport.ValidatorAddress,
			Fees:             msg.BaseReport.Fees,
		},
		ResultsWithTestCase: msg.GetResultsWithTestCase(),
	}, nil
}

// ValidateReport validates if the report is valid to get rewards
func (k Keeper) ValidateReport(ctx sdk.Context, val sdk.ValAddress, req *airequest.AIRequest) error {
	// Check if the validator is in the requested list of validators
	if !k.ContainsVal(req.GetValidators(), val) {
		return sdkerrors.Wrap(types.ErrValidatorNotFound, fmt.Sprintln("failed to find the requested validator"))
	}
	return nil
}
