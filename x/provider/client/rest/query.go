package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oraichain/orai/x/provider/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// TODO: Define your GET REST AIDataSources
	r.HandleFunc(
		fmt.Sprintf("/%s/oscript/{%s}", storeName, restName),
		queryOracleScriptHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/onames", storeName),
		queryOracleScriptNamesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/oscripts", storeName),
		queryOracleScriptsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/datasource/{%s}", storeName, restName),
		queryDataSourceHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/datasources", storeName),
		queryDataSourcesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/dnames", storeName),
		queryDataSourceNamesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/testcase/{%s}", storeName, restName),
		queryTestCaseHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/testcases", storeName),
		queryTestCasesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/tcnames", storeName),
		queryTestCaseNamesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/min_fees/{%s}", storeName, restName),
		queryMinimumFeesHandlerFn(cliCtx),
	).Methods("GET")
}

func queryOracleScriptHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/oscript/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryOracleScriptsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get query params
		v := r.URL.Query()
		page := v.Get("page")
		limit := v.Get("limit")
		name := v.Get("name")

		// In case the request does not include pagination parameters
		if page == "" || limit == "" {
			page = types.DefaultQueryPage
			limit = types.DefaultQueryLimit
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/oscripts", storeName), []byte(page+"-"+limit+"-"+name))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryOracleScriptNamesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/onames", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDataSourceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/datasource/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDataSourcesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		page := v.Get("page")
		limit := v.Get("limit")
		name := v.Get("name")

		// In case the request does not include pagination parameters
		if page == "" || limit == "" {
			page = types.DefaultQueryPage
			limit = types.DefaultQueryLimit
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/datasources", storeName), []byte(page+"-"+limit+"-"+name))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDataSourceNamesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/dnames", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTestCaseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/testcase/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTestCasesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()
		page := v.Get("page")
		limit := v.Get("limit")
		name := v.Get("name")

		// In case the request does not include pagination parameters
		if page == "" || limit == "" {
			page = types.DefaultQueryPage
			limit = types.DefaultQueryLimit
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/testcases", storeName), []byte(page+"-"+limit+"-"+name))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTestCaseNamesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tcnames", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryMinimumFeesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]
		v := r.URL.Query()
		valNum := v.Get("val_num")

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/min_fees/%s", storeName, name), []byte(valNum))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
