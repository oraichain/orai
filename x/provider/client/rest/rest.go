package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	storeName = "provider"
	restName  = storeName
)

// RegisterRoutes registers provider-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
	registerTxRoutes(clientCtx, r)
}
