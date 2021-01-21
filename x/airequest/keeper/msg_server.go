package keeper

import (
	"context"
	"fmt"
	"strings"

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
	validators, err := k.keeper.RandomValidators(ctx, int(msg.ValidatorCount), []byte(msg.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoinsNormalized(msg.Fees)
	// Compute the fee allocated for oracle module to distribute to active validators.
	rewardRatio := sdk.NewDecWithPrec(k.keeper.ProviderKeeper.GetKeyOracleScriptRewardPercentage(ctx), 2)
	// We need to calculate the final 70% fee given by the user because the remaining 30% must be reserved for the proposer and validators.
	providedCoins, _ := sdk.NewDecCoinsFromCoins(fees...).MulDecTruncate(rewardRatio).TruncateDecimal()

	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.keeper.ProviderKeeper.GetDNamesTcNames(ctx, msg.OracleScriptName)
	if err != nil {
		return nil, err
	}

	// collect data source and test case objects to store into the request
	dataSourceObjs, testcaseObjs, err := k.getDSourcesTCases(ctx, aiDataSources, testCases)
	if err != nil {
		return nil, err
	}

	finalFees, err := k.keeper.ProviderKeeper.GetMinimumFees(ctx, aiDataSources, testCases, len(validators))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Error getting minimum fees from oracle script")
	}
	fmt.Println("final fees needed: ", finalFees.String())

	// If the total fee is larger than the fee provided by the user then we return error
	if finalFees.IsAnyGT(providedCoins) {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Fees given by the users are less than the total fees needed")
	}
	// set a new request with the aggregated result into blockchain
	request := types.NewAIRequest(msg.RequestID, msg.OracleScriptName, msg.Creator, validators, ctx.BlockHeight(), dataSourceObjs, testcaseObjs, fees, msg.Input, msg.ExpectedOutput)

	k.keeper.SetAIRequest(ctx, request.RequestID, request)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeSetAIRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
	)
	for _, validator := range validators {
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeRequestValidator, validator.String()),
		)
	}
	ctx.EventManager().EmitEvent(event)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event = sdk.NewEvent(types.EventTypeRequestWithData)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
		sdk.NewAttribute(types.AttributeOracleScriptName, request.OracleScriptName),
		sdk.NewAttribute(types.AttributeRequestCreator, msg.Creator.String()),
		sdk.NewAttribute(types.AttributeRequestValidatorCount, fmt.Sprint(msg.ValidatorCount)),
		sdk.NewAttribute(types.AttributeRequestInput, string(msg.Input)),
		sdk.NewAttribute(types.AttributeRequestExpectedOutput, string(msg.ExpectedOutput)),
		sdk.NewAttribute(types.AttributeRequestDSources, strings.Join(aiDataSources, "-")),
		sdk.NewAttribute(types.AttributeRequestTCases, strings.Join(testCases, "-")),
	)
	ctx.EventManager().EmitEvent(event)

	return types.NewMsgSetAIRequestRes(request.GetRequestID(), request.GetOracleScriptName(), request.GetCreator(), request.GetFees().String(), msg.GetValidatorCount(), request.GetInput(), request.GetExpectedOutput()), nil
}

func (k msgServer) getDSourcesTCases(ctx sdk.Context, dSources, tCases []string) (dSourceObjs []provider.AIDataSource, tCaseObjs []provider.TestCase, errors error) {

	// collect data source objects
	for _, dSource := range dSources {
		dSourceObj, err := k.keeper.ProviderKeeper.GetAIDataSource(ctx, dSource)
		if err != nil {
			return nil, nil, err
		}
		dSourceObjs = append(dSourceObjs, *dSourceObj)
	}

	// collect test case objects
	for _, tCase := range tCases {
		tCaseObj, err := k.keeper.ProviderKeeper.GetTestCase(ctx, tCase)
		if err != nil {
			return nil, nil, err
		}
		tCaseObjs = append(tCaseObjs, *tCaseObj)
	}
	return dSourceObjs, tCaseObjs, nil
}
