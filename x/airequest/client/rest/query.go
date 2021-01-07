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
		fmt.Sprintf("/%s/aireq/{%s}", storeName, restName),
		queryAIRequestHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireqs", storeName),
		queryAIRequestIDsHandlerFn(cliCtx),
	).Methods("GET")
}

func queryAIRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireq/%s", storeName, id), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryAIRequestIDsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireqs", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
