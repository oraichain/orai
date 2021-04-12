package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/aioracle/types"
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

func (k msgServer) CreateAIOracle(goCtx context.Context, msg *types.MsgSetAIOracleReq) (*types.MsgSetAIOracleRes, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators, err := k.keeper.RandomValidators(ctx, int(msg.ValidatorCount), []byte(msg.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	providedFees, _ := sdk.ParseCoinsNormalized(msg.Fees)

	// check if the account has enough spendable coins
	spendableCoins := k.keeper.BankKeeper.SpendableCoins(ctx, msg.Creator)
	// If the total fee is larger or equal to the spendable coins of the user then we return error
	if providedFees.IsAnyGTE(spendableCoins) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Your account has run out of tokens to create the AI Request\n"))
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Your account has run out of tokens to create the AI Request")
	}

	// substract coins in the creator wallet to charge fees
	err = k.keeper.BankKeeper.SubtractCoins(ctx, msg.Creator, providedFees)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Your account has run out of tokens to create the AI Request, or there is something wrong")
	}

	// set a new request with the aggregated result into blockchain
	request := types.NewAIOracle(msg.RequestID, msg.Contract, msg.Creator, validators, ctx.BlockHeight(), providedFees, msg.Input)

	k.keeper.SetAIOracle(ctx, request.RequestID, request)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequestWithData)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
		sdk.NewAttribute(types.AttributeContract, request.Contract.String()),
		sdk.NewAttribute(types.AttributeRequestCreator, msg.Creator.String()),
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
		request.GetRequestID(), request.GetContract(),
		request.GetCreator(), request.GetFees().String(), msg.GetValidatorCount(),
		request.GetInput(),
	), nil
}
