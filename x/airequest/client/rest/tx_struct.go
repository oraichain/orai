package rest

//go:generate msgp

type SetAIRequestReq struct {
	From             string `json:"from" msg:"from"`
	Memo             string `json:"memo" msg:"memo"`
	ChainID          string `json:"chain_id" msg:"chain_id"`
	AccountNumber    uint64 `json:"account_number" msg:"account_number"`
	Sequence         uint64 `json:"sequence" msg:"sequence"`
	GasPrices        string `json:"gas_prices" msg:"gas_prices"`
	Gas              string `json:"gas" msg:"gas"`
	GasAdjustment    string `json:"gas_adjustment" msg:"gas_adjustment"`
	Simulate         bool   `json:"simulate" msg:"simulate"`
	OracleScriptName string `json:"oracle_script_name" msg:"oracle_script_name"`
	Input            []byte `json:"input" msg:"input"`
	ExpectedOutput   []byte `json:"expected_output" msg:"expected_output"`
	Fees             string `json:"fees" msg:"fees"`
	ValidatorCount   int    `json:"validator_count" msg:"validator_count"`
}
