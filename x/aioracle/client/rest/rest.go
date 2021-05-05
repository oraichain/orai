package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	storeName = "aioracle"
	restName  = storeName
)

// RegisterRoutes registers aioracle-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {

}
