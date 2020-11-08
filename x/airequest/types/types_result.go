package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airequest/exported"
)

// implements AIRequestResultI interface
var _ exported.AIRequestResultI = AIRequestResult{}

// implements ValResultI interface
var _ exported.ValResultI = ValResult{}

// AIRequestResult stores the final result after aggregating the results from the reports of an AI request
type AIRequestResult struct {
	RequestID string     `json:"request_id"`
	Results   ValResults `json:"results"`
	Status    string     `json:"request_status"`
}

// ValResult stores the result information from a validator that has executed the oracle script
type ValResult struct {
	Validator sdk.ValAddress `json:"validator_address"`
	Result    []byte         `json:"result"`
}

// ValResults is the list of results struct
type ValResults []ValResult

// NewAIRequestResult is a constructor for the ai request result struct
func NewAIRequestResult(
	requestID string,
	results ValResults,
	status string,
) AIRequestResult {
	return AIRequestResult{
		RequestID: requestID,
		Results:   results,
		Status:    status,
	}
}

// NewValResult is a constructor for the validator result
func NewValResult(
	val sdk.ValAddress,
	result []byte,
) ValResult {
	return ValResult{
		Validator: val,
		Result:    result,
	}
}
