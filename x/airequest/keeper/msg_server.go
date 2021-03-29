package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/types"
	provider "github.com/oraichain/orai/x/provider/types"
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

func (k msgServer) CreateAIRequest(goCtx context.Context, msg *types.MsgSetAIRequest) (*types.MsgSetAIRequestRes, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate if the oracle script exists or not
	if !k.keeper.providerKeeper.HasOracleScript(ctx, msg.OracleScriptName) {
		return nil, types.ErrOScriptNotFound
	}

	// validate if the request id exists or not
	if k.keeper.HasAIRequest(ctx, msg.RequestID) {
		return nil, sdkerrors.Wrap(types.ErrRequestInvalid, "The request id already exists")
	}

	validators, err := k.keeper.RandomValidators(ctx, int(msg.ValidatorCount), []byte(msg.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	providedFees, _ := sdk.ParseCoinsNormalized(msg.Fees)

	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.keeper.providerKeeper.GetDNamesTcNames(ctx, msg.OracleScriptName)
	if err != nil {
		return nil, err
	}

	// collect data source and test case objects to store into the request
	dataSourceObjs, testcaseObjs, err := k.getDSourcesTCases(ctx, aiDataSources, testCases)
	if err != nil {
		return nil, err
	}

	requiredFees, err := k.keeper.providerKeeper.GetMinimumFees(ctx, aiDataSources, testCases, len(validators), k.keeper.providerKeeper.GetOracleScriptRewardPercentageParam(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Error getting minimum fees from oracle script")
	}
	k.keeper.Logger(ctx).Info(fmt.Sprintf("required fees needed: %v\n", requiredFees))

	// If the total fee is larger than the fee provided by the user then we return error
	if requiredFees.IsAnyGT(providedFees) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Your account has run out of tokens to create the AI Request\n"))
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Fees given by the users are less than the total fees needed")
	}

	// check if the account has enough spendable coins
	spendableCoins := k.keeper.bankKeeper.SpendableCoins(ctx, msg.Creator)
	// If the total fee is larger or equal to the spendable coins of the user then we return error
	if requiredFees.IsAnyGTE(spendableCoins) || providedFees.IsAnyGTE(spendableCoins) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Your account has run out of tokens to create the AI Request\n"))
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Your account has run out of tokens to create the AI Request")
	}

	// substract coins in the creator wallet to charge fees
	err = k.keeper.bankKeeper.SubtractCoins(ctx, msg.Creator, providedFees)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Your account has run out of tokens to create the AI Request, or there is something wrong")
	}

	// set a new request with the aggregated result into blockchain
	request := types.NewAIRequest(msg.RequestID, msg.OracleScriptName, msg.Creator, validators, ctx.BlockHeight(), dataSourceObjs, testcaseObjs, providedFees, msg.Input, msg.ExpectedOutput)

	k.keeper.SetAIRequest(ctx, request.RequestID, request)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequestWithData)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
		sdk.NewAttribute(types.AttributeOracleScriptName, request.OracleScriptName),
		sdk.NewAttribute(types.AttributeRequestCreator, msg.Creator.String()),
		sdk.NewAttribute(types.AttributeRequestValidatorCount, fmt.Sprint(msg.ValidatorCount)),
		sdk.NewAttribute(types.AttributeRequestInput, string(msg.Input)),
		sdk.NewAttribute(types.AttributeRequestExpectedOutput, string(msg.ExpectedOutput)),
	)

	for _, validator := range validators {
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeRequestValidator, validator.String()),
		)
	}

	// these are multiple attribute for array
	for _, aiDataSource := range aiDataSources {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeRequestDSources, aiDataSource))
	}
	for _, testCase := range testCases {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeRequestTCases, testCase))
	}

	ctx.EventManager().EmitEvent(event)

	return types.NewMsgSetAIRequestRes(
		request.GetRequestID(), request.GetOracleScriptName(),
		request.GetCreator(), request.GetFees().String(), msg.GetValidatorCount(),
		request.GetInput(), request.GetExpectedOutput(),
	), nil
}

func (k msgServer) getDSourcesTCases(ctx sdk.Context, dSources, tCases []string) (dSourceObjs []provider.AIDataSource, tCaseObjs []provider.TestCase, errors error) {

	// collect data source objects
	for _, dSource := range dSources {
		dSourceObj, err := k.keeper.providerKeeper.GetAIDataSource(ctx, dSource)
		if err != nil {
			return nil, nil, err
		}
		dSourceObjs = append(dSourceObjs, *dSourceObj)
	}

	// collect test case objects
	for _, tCase := range tCases {
		tCaseObj, err := k.keeper.providerKeeper.GetTestCase(ctx, tCase)
		if err != nil {
			return nil, nil, err
		}
		tCaseObjs = append(tCaseObjs, *tCaseObj)
	}
	return dSourceObjs, tCaseObjs, nil
}
