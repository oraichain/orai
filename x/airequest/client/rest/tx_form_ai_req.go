package rest

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type setFormAIRequestReq struct {
	OracleScriptName string       `json:"oracle_script_name"`
	From             string       `json:"from"`
	Memo             string       `json:"memo"`
	ChainID          string       `json:"chain_id"`
	AccountNumber    uint64       `json:"account_number"`
	Sequence         uint64       `json:"sequence"`
	Fees             sdk.Coins    `json:"fees"`
	GasPrices        sdk.DecCoins `json:"gas_prices"`
	Gas              string       `json:"gas"`
	GasAdjustment    string       `json:"gas_adjustment"`
	Simulate         bool         `json:"simulate"`
	Input            string       `json:"input"`
	ExpectedOutput   string       `json:"expected_output"`
	ValidatorCount   int          `json:"validator_count"`
}

// TODO: Need GETTER SETTER here

func setAIRequestHandlerFn(cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request) setFormAIRequestReq {
	// convert from form-data string to correct data type of baseReq
	accNum, err := strconv.Atoi(r.FormValue("account_number"))
	if len(r.FormValue("account_number")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for account number")
		return setFormAIRequestReq{}
	}

	// convert from form-data string to correct data type of baseReq
	sequence, err := strconv.Atoi(r.FormValue("sequence"))
	if len(r.FormValue("sequence")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for sequence")
		return setFormAIRequestReq{}
	}

	// convert from form-data string to correct data type of baseReq
	simulate, err := strconv.ParseBool(r.FormValue("simulate"))
	if len(r.FormValue("simulate")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for simulate")
		return setFormAIRequestReq{}
	}

	// convert from form-data string to correct data type of baseReq
	feeStr := r.FormValue("fees")
	var fees sdk.Coins

	if len(feeStr) != 0 {
		fees, err = sdk.ParseCoins(feeStr)
		if len(r.FormValue("fees")) != 0 && err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for fees")
			return setFormAIRequestReq{}
		}
	}

	// convert from form-data string to correct data type of baseReq
	gasPricesStr := r.FormValue("gas_prices")
	var gasPrices sdk.DecCoins

	if len(gasPricesStr) != 0 {
		gasPrices, err = sdk.ParseDecCoins(gasPricesStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for gas price")
			return setFormAIRequestReq{}
		}
	}

	// convert from form-data string to correct data type of baseReq
	valCount, err := strconv.Atoi(r.FormValue("validator_count"))
	if len(r.FormValue("validator_count")) != 0 && err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request for validator count")
		return setFormAIRequestReq{}
	}

	// This may be redundant, but it does not affect much so leave this for now.
	req := setFormAIRequestReq{
		OracleScriptName: r.FormValue("oracle_script_name"),
		From:             r.FormValue("from"),
		Memo:             r.FormValue("memo"),
		ChainID:          r.FormValue("chain_id"),
		AccountNumber:    uint64(accNum),
		Sequence:         uint64(sequence),
		Fees:             fees,
		GasPrices:        gasPrices,
		Gas:              r.FormValue("gas"),
		GasAdjustment:    r.FormValue("gas_adjustment"),
		Simulate:         simulate,
		Input:            r.FormValue("input"),
		ExpectedOutput:   r.FormValue("expected_output"),
		ValidatorCount:   valCount,
	}

	return req
}
