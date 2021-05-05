package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/aioracle/types"
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

	requiredFees, err := k.querier.calculateMinimumFees(ctx, msg.TestOnly, contract, len(validators))
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

	return types.NewMsgSetAIOracleRes(
		request.GetRequestID(), msg.GetContract(),
		msg.GetCreator(), request.GetFees().String(), msg.GetValidatorCount(),
		request.GetInput(), request.GetTestOnly(),
	), nil
}
