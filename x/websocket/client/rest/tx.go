package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (

	// "bytes"
	// "net/http"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	// "github.com/cosmos/cosmos-sdk/types/rest"
	// "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type setPriceRequestReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	OracleScriptName string       `json:"oracle_script_name"`
	Input            string       `json:"input"`
	ExpectedOutput   string       `json:"expected_output"`
	Fees             string       `json:"fees"`
	ValidatorCount   int          `json:"validator_count"`
}

type setKYCRequestReq struct {
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
}

type setIPFSImage struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
}
