package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/ai-request/types"
)

// NewQuerier creates a new querier for provider clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryAIRequest:
			return queryAIRequest(ctx, path[1:], keeper)
		case types.QueryAIRequestIDs:
			return queryAIRequestIDs(ctx, keeper)
		case types.QueryFullRequest:
			return queryFullRequestByID(ctx, path[1:], keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown provider query")
		}
	}
}

// queryAIRequest queries an AI request given its request ID
func queryAIRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	// tsao cho nay lai lay path[0] ?
	request, err := keeper.GetAIRequest(ctx, path[0])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrRequestNotFound, err.Error())
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResAIRequest(request.RequestID, request.Creator, request.OracleScriptName, request.Validators, request.BlockHeight, request.AIDataSources, request.TestCases, request.Fees.String()))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryFullRequestByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "error")
	}
	// id of the request
	id := path[0]
	request, err := k.GetAIRequest(ctx, id)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRequestNotFound, err.Error())
	}
	// collect all the reports of a given request id
	reports := k.GetReports(ctx, id)

	// collect the result of a given request id

	result, err := k.GetResult(ctx, id)
	if err != nil {
		// init a nil result showing that the result does not have a result yet
		result = types.NewAIRequestResult(id, nil, types.RequestStatusPending)
	}

	res, err := codec.MarshalJSONIndent(k.cdc, types.NewQueryResFullRequest(request, reports, result))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryAIRequestIDs returns all the AI Request IDs in the store
func queryAIRequestIDs(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var requestIDs types.QueryResAIRequestIDs

	iterator := keeper.GetAllAIRequestIDs(ctx)

	for ; iterator.Valid(); iterator.Next() {
		requestIDs = append(requestIDs, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, requestIDs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
