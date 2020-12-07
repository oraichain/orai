package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// TODO: Define your GET REST AIDataSources

	r.HandleFunc(
		fmt.Sprintf("/%s/fullreq/{%s}", storeName, restName),
		queryFullAIRequestHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/reward/{%s}", storeName, restName),
		queryRewardHandlerFn(cliCtx),
	).Methods("GET")
}

func queryFullAIRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/fullreq/%s", storeName, id), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRewardHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/reward/%s", storeName, id), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
