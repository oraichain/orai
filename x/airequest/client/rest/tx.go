package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (

	// "bytes"
	// "net/http"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/segmentio/ksuid"
	"github.com/tinylib/msgp/msgp"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	// "github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq", storeName),
		setAIRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/%s/aireq/testreq", storeName),
		setTestRequestHandlerFn(cliCtx),
	).Methods("POST")
}

type setTestRequestReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	OracleScriptName string       `json:"oracle_script_name"`
	Input            []byte       `json:"input"`
	ExpectedOutput   []byte       `json:"expected_output"`
	Fees             string       `json:"fees"`
	ValidatorCount   int          `json:"validator_count"`
}

func setTestRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req setTestRequestReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "AVXSD")
			return
		}

		// create the message
		msg := types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees, req.ValidatorCount, req.Input, req.ExpectedOutput)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GHYK")
			return
		}
		// Collect fees in Coins type. Bug: cannot set fee through json using REST API => This is the workaround
		fees, _ := sdk.ParseCoins(req.Fees)
		baseReq.Fees = fees
		if !baseReq.ValidateBasic(w) {
			return
		}
		fmt.Println("base req: ", baseReq)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func setAIRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req SetAIRequestReq

		if !ReadRESTReq(w, r, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		fmt.Println("read req after parsing: ", req)
		// Collect fees & gas prices in Coins type. Bug: cannot set fee through json using REST API => This is the workaround
		fees, _ := sdk.ParseCoins(req.Fees)
		gas, _ := sdk.ParseCoins(req.GasPrices)
		baseReq := rest.NewBaseReq(req.From, req.Memo, req.ChainID, req.Gas, req.GasAdjustment, req.AccountNumber, req.Sequence, fees, sdk.NewDecCoinsFromCoins(gas...), req.Simulate)
		baseReq.Fees = fees
		if !baseReq.ValidateBasic(w) {
			return
		}
		fmt.Println("base req: ", baseReq)
		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "AVXSD")
			return
		}

		// create the message
		msg := types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees, req.ValidatorCount, req.Input, req.ExpectedOutput)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GHYK")
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

// ReadRESTReq reads and unmarshals a Request's body to the the BaseReq stuct.
// Writes an error response to ResponseWriter and returns true if errors occurred.
func ReadRESTReq(w http.ResponseWriter, r *http.Request, req *SetAIRequestReq) bool {

	err := msgp.Decode(r.Body, req)
	if err != nil {
		fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBB: ", err)
		rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to decode JSON payload: %s", err))
		return false
	}

	return true
}
