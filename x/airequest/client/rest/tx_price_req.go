package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/segmentio/ksuid"
)

type setPriceRequestReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	OracleScriptName string       `json:"oracle_script_name"`
	Input            string       `json:"input"`
	ExpectedOutput   string       `json:"expected_output"`
	Fees             string       `json:"fees"`
	ValidatorCount   int          `json:"validator_count"`
}

func setPriceRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req setPriceRequestReq

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
		msg := types.NewMsgSetPriceRequest(types.NewMsgSetAIRequest(ksuid.New().String(), req.OracleScriptName, addr, req.Fees, req.ValidatorCount, req.Input, req.ExpectedOutput))
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
