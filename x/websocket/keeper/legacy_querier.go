package keeper

import (
	// this line is used by starport scaffolding # 1

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/websocket/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

// Query statements for data sources
const (
// TODO: Describe query parameters, update <action> with your query
// Query<Action>    = "<action>"
)

// NewLegacyQuerier implementing for legacy querier
func NewLegacyQuerier(querier *Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		// var (
		// 	rsp interface{}
		// 	err error
		// )
		// goCtx := ctx.Context()

		// switch path[0] {

		// default:
		// 	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		// }

		// // wrap serialization
		// if err != nil {
		// 	return nil, err
		// }
		// bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, rsp)
		// if err != nil {
		// 	return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		// }
		// return bz, nil
	}
}
