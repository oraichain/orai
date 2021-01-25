package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (

	// "bytes"
	// "net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	// "github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

// type createOracleScriptReq struct {
// 	BaseReq     rest.BaseReq `json:"base_req"`
// 	Name        string       `json:"name"`
// 	Contract    string       `json:"contract"`
// 	Fees        string       `json:"fees"`
// 	Description string       `json:"description"`
// 	DataSources []string     `json:"data_sources"`
// 	TestCases   []string     `json:"test_cases"`
// }

// type editOracleScriptReq struct {
// 	BaseReq     rest.BaseReq `json:"base_req"`
// 	OldName     string       `json:"old_name"`
// 	NewName     string       `json:"new_name"`
// 	Contract    string       `json:"contract"`
// 	Fees        string       `json:"fees"`
// 	Description string       `json:"description"`
// 	DataSources []string     `json:"data_sources"`
// 	TestCases   []string     `json:"test_cases"`
// }

// type createDataSourceReq struct {
// 	BaseReq     rest.BaseReq `json:"base_req"`
// 	Name        string       `json:"name"`
// 	Contract    string       `json:"contract"`
// 	Fees        string       `json:"fees"`
// 	Description string       `json:"description"`
// }

// type editDataSourceReq struct {
// 	BaseReq     rest.BaseReq `json:"base_req"`
// 	OldName     string       `json:"old_name"`
// 	NewName     string       `json:"new_name"`
// 	Contract    string       `json:"contract"`
// 	Fees        string       `json:"fees"`
// 	Description string       `json:"description"`
// }

// type createTestCaseReq struct {
// 	BaseReq     rest.BaseReq `json:"base_req"`
// 	Name        string       `json:"name"`
// 	Contract    string       `json:"contract"`
// 	Fees        string       `json:"fees"`
// 	Description string       `json:"description"`
// }

// type editTestCaseReq struct {
// 	BaseReq     rest.BaseReq `json:"base_req"`
// 	OldName     string       `json:"old_name"`
// 	NewName     string       `json:"new_name"`
// 	Contract    string       `json:"contract"`
// 	Fees        string       `json:"fees"`
// 	Description string       `json:"description"`
// }

func registerTxRoutes(clientCtx client.Context, r *mux.Router) {
	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/oscript", storeName),
	// 	setOracleScriptHandlerFn(clientCtx),
	// ).Methods("POST")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/oscript", storeName),
	// 	editOracleScriptHandlerFn(clientCtx),
	// ).Methods("PATCH")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/datasource", storeName),
	// 	setDataSourceHandlerFn(clientCtx),
	// ).Methods("POST")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/datasource", storeName),
	// 	editDataSourceHandlerFn(clientCtx),
	// ).Methods("PATCH")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/testcase", storeName),
	// 	setTestCaseHandlerFn(clientCtx),
	// ).Methods("POST")

	// r.HandleFunc(
	// 	fmt.Sprintf("/%s/testcase", storeName),
	// 	editTestCaseHandlerFn(clientCtx),
	// ).Methods("PATCH")
}

// func setOracleScriptHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req createOracleScriptReq

// 		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
// 			return
// 		}

// 		baseReq := req.BaseReq.Sanitize()

// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		// collect valid address from the request address string
// 		addr, err := sdk.AccAddressFromBech32(baseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid account address owner")
// 			return
// 		}

// 		// create the message
// 		msg := types.NewMsgCreateOracleScript(req.Name, req.Contract, addr, req.Fees, req.Description, req.DataSources, req.TestCases)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		tx.WriteGeneratedTxResponse(clientCtx, w, baseReq, msg)
// 	}
// }

// func editOracleScriptHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req editOracleScriptReq

// 		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
// 			return
// 		}

// 		baseReq := req.BaseReq.Sanitize()

// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		// collect valid address from the request address string
// 		addr, err := sdk.AccAddressFromBech32(baseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		// create the message
// 		msg := types.NewMsgEditOracleScript(req.OldName, req.NewName, req.Contract, addr, req.Fees, req.Description, req.DataSources, req.TestCases)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		tx.WriteGeneratedTxResponse(clientCtx, w, baseReq, msg)
// 	}
// }

// func setDataSourceHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req createDataSourceReq

// 		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
// 			return
// 		}

// 		baseReq := req.BaseReq.Sanitize()

// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		// collect valid address from the request address string
// 		addr, err := sdk.AccAddressFromBech32(baseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid account address owner")
// 			return
// 		}

// 		// create the message
// 		msg := types.NewMsgCreateAIDataSource(req.Name, req.Contract, addr, req.Fees, req.Description)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		tx.WriteGeneratedTxResponse(clientCtx, w, baseReq, msg)
// 	}
// }

// func editDataSourceHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req editDataSourceReq

// 		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
// 			return
// 		}

// 		baseReq := req.BaseReq.Sanitize()

// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		// collect valid address from the request address string
// 		addr, err := sdk.AccAddressFromBech32(baseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		// create the message
// 		msg := types.NewMsgEditAIDataSource(req.OldName, req.NewName, req.Contract, addr, req.Fees, req.Description)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		tx.WriteGeneratedTxResponse(clientCtx, w, baseReq, msg)
// 	}
// }

// func setTestCaseHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req createTestCaseReq

// 		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
// 			return
// 		}

// 		baseReq := req.BaseReq.Sanitize()

// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		// collect valid address from the request address string
// 		addr, err := sdk.AccAddressFromBech32(baseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid account address owner")
// 			return
// 		}

// 		// create the message
// 		msg := types.NewMsgCreateTestCase(req.Name, req.Contract, addr, req.Fees, req.Description)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		tx.WriteGeneratedTxResponse(clientCtx, w, baseReq, msg)
// 	}
// }

// func editTestCaseHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req editTestCaseReq

// 		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
// 			return
// 		}

// 		baseReq := req.BaseReq.Sanitize()

// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		// collect valid address from the request address string
// 		addr, err := sdk.AccAddressFromBech32(baseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		// create the message
// 		msg := types.NewMsgEditTestCase(req.OldName, req.NewName, req.Contract, addr, req.Fees, req.Description)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		tx.WriteGeneratedTxResponse(clientCtx, w, baseReq, msg)
// 	}
// }
