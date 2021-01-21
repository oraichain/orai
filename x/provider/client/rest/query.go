package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// TODO: Define your GET REST AIDataSources
	r.HandleFunc(
		fmt.Sprintf("/%s/oscript/{%s}", storeName, restName),
		queryOracleScriptHandlerFn(clientCtx),
	).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/onames", storeName),
	// 	queryOracleScriptNamesHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/oscripts", storeName),
	// 	queryOracleScriptsHandlerFn(clientCtx),
	// ).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/datasource/{%s}", storeName, restName),
		queryDataSourceHandlerFn(clientCtx),
	).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/datasources", storeName),
	// 	queryDataSourcesHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/dnames", storeName),
	// 	queryDataSourceNamesHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/testcase/{%s}", storeName, restName),
	// 	queryTestCaseHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/testcases", storeName),
	// 	queryTestCasesHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/tcnames", storeName),
	// 	queryTestCaseNamesHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/min_fees/{%s}", storeName, restName),
	// 	queryMinimumFeesHandlerFn(clientCtx),
	// ).Methods("GET")
}

func queryOracleScriptHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/oscript/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

// func queryOracleScriptsHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Get query params
// 		v := r.URL.Query()
// 		page := v.Get("page")
// 		limit := v.Get("limit")
// 		name := v.Get("name")

// 		// In case the request does not include pagination parameters
// 		if page == "" || limit == "" {
// 			page = types.DefaultQueryPage
// 			limit = types.DefaultQueryLimit
// 		}

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/oscripts", storeName), []byte(page+"-"+limit+"-"+name))
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryOracleScriptNamesHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/onames", storeName), nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

func queryDataSourceHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/datasource/%s", storeName, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

// func queryDataSourcesHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		v := r.URL.Query()
// 		page := v.Get("page")
// 		limit := v.Get("limit")
// 		name := v.Get("name")

// 		// In case the request does not include pagination parameters
// 		if page == "" || limit == "" {
// 			page = types.DefaultQueryPage
// 			limit = types.DefaultQueryLimit
// 		}

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/datasources", storeName), []byte(page+"-"+limit+"-"+name))
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryDataSourceNamesHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/dnames", storeName), nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryTestCaseHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		name := vars[restName]

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/testcase/%s", storeName, name), nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryTestCasesHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		v := r.URL.Query()
// 		page := v.Get("page")
// 		limit := v.Get("limit")
// 		name := v.Get("name")

// 		// In case the request does not include pagination parameters
// 		if page == "" || limit == "" {
// 			page = types.DefaultQueryPage
// 			limit = types.DefaultQueryLimit
// 		}

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/testcases", storeName), []byte(page+"-"+limit+"-"+name))
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryTestCaseNamesHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/tcnames", storeName), nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryMinimumFeesHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		name := vars[restName]
// 		v := r.URL.Query()
// 		valNum := v.Get("val_num")

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/min_fees/%s", storeName, name), []byte(valNum))
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }
