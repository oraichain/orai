package airequest

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
	provider "github.com/oraichain/orai/x/provider/exported"
)

// NewHandler creates an sdk.Handler for all the airequest type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetAIRequest:
			return handleMsgSetAIRequest(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgSetAIRequest(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetAIRequest) (*sdk.Result, error) {
	validators, err := k.RandomValidators(ctx, msg.ValidatorCount, []byte(msg.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.Fees)
	// Compute the fee allocated for oracle module to distribute to active validators.
	rewardRatio := sdk.NewDecWithPrec(k.ProviderKeeper.GetKeyOracleScriptRewardPercentage(ctx), 2)
	// We need to calculate the final 70% fee given by the user because the remaining 30% must be reserved for the proposer and validators.
	providedCoins, _ := sdk.NewDecCoinsFromCoins(fees...).MulDecTruncate(rewardRatio).TruncateDecimal()

	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.ProviderKeeper.GetDNamesTcNames(ctx, msg.OracleScriptName)
	if err != nil {
		return nil, err
	}

	// collect data source and test case objects to store into the request
	dataSourceObjs, testcaseObjs, err := getDSourcesTCases(ctx, k, aiDataSources, testCases)
	if err != nil {
		return nil, err
	}

	finalFees, err := k.ProviderKeeper.GetMinimumFees(ctx, aiDataSources, testCases, len(validators))
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

	k.SetAIRequest(ctx, request.RequestID, request)

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
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func getDSourcesTCases(ctx sdk.Context, k Keeper, dSources, tCases []string) (dSourceObjs []provider.AIDataSourceI, tCaseObjs []provider.TestCaseI, errors error) {

	// collect data source objects
	for _, dSource := range dSources {
		dSourceObj, err := k.ProviderKeeper.GetAIDataSourceI(ctx, dSource)
		if err != nil {
			return nil, nil, err
		}
		dSourceObjs = append(dSourceObjs, dSourceObj)
	}

	// collect test case objects
	for _, tCase := range tCases {
		tCaseObj, err := k.ProviderKeeper.GetTestCaseI(ctx, tCase)
		if err != nil {
			return nil, nil, err
		}
		tCaseObjs = append(tCaseObjs, tCaseObj)
	}
	return dSourceObjs, tCaseObjs, nil
}
