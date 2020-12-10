package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (

	// "bytes"
	// "net/http"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oraichain/orai/x/provider/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	// "github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type createOracleScriptReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Name        string       `json:"name"`
	CodePath    string       `json:"code_path"`
	Description string       `json:"description"`
	DataSources []string     `json:"data_sources"`
	TestCases   []string     `json:"test_cases"`
}

type editOracleScriptReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	OldName     string       `json:"old_name"`
	NewName     string       `json:"new_name"`
	CodePath    string       `json:"code_path"`
	Description string       `json:"description"`
	DataSources []string     `json:"data_sources"`
	TestCases   []string     `json:"test_cases"`
}

type createDataSourceReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Name        string       `json:"name"`
	CodePath    string       `json:"code_path"`
	Fees        string       `json:"fees"`
	Description string       `json:"description"`
}

type editDataSourceReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	OldName     string       `json:"old_name"`
	NewName     string       `json:"new_name"`
	CodePath    string       `json:"code_path"`
	Fees        string       `json:"fees"`
	Description string       `json:"description"`
}

type createTestCaseReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Name        string       `json:"name"`
	CodePath    string       `json:"code_path"`
	Fees        string       `json:"fees"`
	Description string       `json:"description"`
}

type editTestCaseReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	OldName     string       `json:"old_name"`
	NewName     string       `json:"new_name"`
	CodePath    string       `json:"code_path"`
	Fees        string       `json:"fees"`
	Description string       `json:"description"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/%s/oscript", storeName),
		setOracleScriptHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/oscript", storeName),
		editOracleScriptHandlerFn(cliCtx),
	).Methods("PATCH")

	r.HandleFunc(
		fmt.Sprintf("/%s/datasource", storeName),
		setDataSourceHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/datasource", storeName),
		editDataSourceHandlerFn(cliCtx),
	).Methods("PATCH")

	r.HandleFunc(
		fmt.Sprintf("/%s/testcase", storeName),
		setTestCaseHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/testcase", storeName),
		editTestCaseHandlerFn(cliCtx),
	).Methods("PATCH")
}

func setOracleScriptHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createOracleScriptReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid account address owner")
			return
		}

		execBytes, err := getFileBytes(req.CodePath)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintln("cannot read file from the given file path"+err.Error()))
			return
		}

		// create the message
		msg := types.NewMsgCreateOracleScript(req.Name, execBytes, addr, req.Description, req.DataSources, req.TestCases)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func editOracleScriptHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editOracleScriptReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// collect the byte code of the source code based on the path
		execBytes, err := ioutil.ReadFile(req.CodePath)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgEditOracleScript(req.OldName, req.NewName, execBytes, addr, req.Description, req.DataSources, req.TestCases)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func setDataSourceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createDataSourceReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid account address owner")
			return
		}

		execBytes, err := getFileBytes(req.CodePath)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintln("cannot read file from the given file path"+err.Error()))
			return
		}

		// create the message
		msg := types.NewMsgCreateAIDataSource(req.Name, execBytes, addr, req.Fees, req.Description)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func editDataSourceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editDataSourceReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// collect the byte code of the source code based on the path
		execBytes, err := ioutil.ReadFile(req.CodePath)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgEditAIDataSource(req.OldName, req.NewName, execBytes, addr, req.Fees, req.Description)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func setTestCaseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTestCaseReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid account address owner")
			return
		}

		execBytes, err := getFileBytes(req.CodePath)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintln("cannot read file from the given file path"+err.Error()))
			return
		}

		// create the message
		msg := types.NewMsgCreateTestCase(req.Name, execBytes, addr, req.Fees, req.Description)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func editTestCaseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editTestCaseReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// collect the byte code of the source code based on the path
		execBytes, err := ioutil.ReadFile(req.CodePath)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgEditTestCase(req.OldName, req.NewName, execBytes, addr, req.Fees, req.Description)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
