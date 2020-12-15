package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aiRequest "github.com/oraichain/orai/x/airequest"
	"github.com/oraichain/orai/x/airesult/types"
)

// NewQuerier creates a new querier for provider clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryFullRequest:
			return queryFullRequestByID(ctx, path[1:], keeper)
		case types.QueryReward:
			return queryReward(ctx, path[1:], keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown provider query")
		}
	}
}

func queryFullRequestByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "error")
	}
	// id of the request
	id := path[0]
	request, err := k.aiRequestKeeper.GetAIRequest(ctx, id)
	if err != nil {
		return nil, sdkerrors.Wrapf(aiRequest.ErrRequestNotFound, err.Error())
	}
	// collect all the reports of a given request id
	reports := k.webSocketKeeper.GetReports(ctx, id)

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
