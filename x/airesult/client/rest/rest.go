package rest

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	storeName = "airesult"
	restName  = storeName
)

// RegisterRoutes registers airequest-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {

	fmt.Println("air request rest")
	// this line is used by starport scaffolding # 2
	//registerQueryRoutes(clientCtx, r)
	//registerTxRoutes(clientCtx, r)
}
