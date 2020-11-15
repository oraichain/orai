package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/types"
)

// NewQuerier creates a new querier for provider clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryAIRequest:
			return queryAIRequest(ctx, path[1:], keeper)
		case types.QueryAIRequestIDs:
			return queryAIRequestIDs(ctx, keeper)
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

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResAIRequest(request.GetRequestID(), request.GetCreator(), request.GetOScriptName(), request.GetValidators(), request.GetBlockHeight(), request.GetAIDataSources(), request.GetTestCases(), request.GetFees().String()))
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
