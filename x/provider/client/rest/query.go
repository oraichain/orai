package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/types/rest"
// 	"github.com/gorilla/mux"
// 	"github.com/oraichain/orai/x/provider/types"
// )

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// TODO: Define your GET REST AIDataSources
	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/oscript/{%s}", storeName, restName),
	// 	queryOracleScriptHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/onames", storeName),
	// 	queryOracleScriptNamesHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/oscripts", storeName),
	// 	queryOracleScriptsHandlerFn(clientCtx),
	// ).Methods("GET")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/datasource/{%s}", storeName, restName),
	// 	queryDataSourceHandlerFn(clientCtx),
	// ).Methods("GET")

	// // r.HandleFunc(
	// // 	fmt.Sprintf("/%s/datasources", storeName),
	// // 	queryDataSourcesHandlerFn(clientCtx),
	// // ).Methods("GET")

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

// func queryOracleScriptHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		name := vars[restName]

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/oscript/%s", storeName, name), nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryOracleScriptsHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Get query params
// 		v := r.URL.Query()
// 		name := v.Get("name")

// 		// set values from query params
// 		var (
// 			page  int
// 			limit int
// 			err   error
// 		)

// 		if page, err = strconv.Atoi(v.Get("page")); err != nil {
// 			page = types.DefaultQueryPage
// 		}

// 		if limit, err = strconv.Atoi(v.Get("limit")); err != nil {
// 			limit = types.DefaultQueryLimit
// 		}

// 		// create params bytes
// 		params := types.ListOracleScriptsReq{Name: name, Limit: int64(limit), Page: int64(page)}
// 		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
// 		if rest.CheckBadRequestError(w, err) {
// 			return
// 		}

// 		// query response
// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/oscripts", storeName), bz)
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

// func queryDataSourceHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		name := vars[restName]

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/datasource/%s", storeName, name), nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryDataSourcesHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		v := r.URL.Query()
// 		name := v.Get("name")

// 		// set values from query params
// 		var (
// 			page  int
// 			limit int
// 			err   error
// 		)

// 		if page, err = strconv.Atoi(v.Get("page")); err != nil {
// 			page = types.DefaultQueryPage
// 		}

// 		if limit, err = strconv.Atoi(v.Get("limit")); err != nil {
// 			limit = types.DefaultQueryLimit
// 		}

// 		// create params bytes
// 		params := types.ListDataSourcesReq{Name: name, Limit: int64(limit), Page: int64(page)}
// 		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
// 		if rest.CheckBadRequestError(w, err) {
// 			return
// 		}

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/datasources", storeName), bz)
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
// 		name := v.Get("name")

// 		// set values from query params
// 		var (
// 			page  int
// 			limit int
// 			err   error
// 		)

// 		if page, err = strconv.Atoi(v.Get("page")); err != nil {
// 			page = types.DefaultQueryPage
// 		}

// 		if limit, err = strconv.Atoi(v.Get("limit")); err != nil {
// 			limit = types.DefaultQueryLimit
// 		}

// 		// create params bytes
// 		params := types.ListTestCasesReq{Name: name, Limit: int64(limit), Page: int64(page)}
// 		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
// 		if rest.CheckBadRequestError(w, err) {
// 			return
// 		}

// 		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/testcases", storeName), bz)
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
