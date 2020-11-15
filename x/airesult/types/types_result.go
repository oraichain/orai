package types

import (
	"github.com/oraichain/orai/x/airesult/exported"
	webSocket "github.com/oraichain/orai/x/websocket/exported"
)

// implements AIRequestResultI interface
var _ exported.AIRequestResultI = AIRequestResult{}

// AIRequestResult stores the final result after aggregating the results from the reports of an AI request
type AIRequestResult struct {
	RequestID string                 `json:"request_id"`
	Results   []webSocket.ValResultI `json:"results"`
	Status    string                 `json:"request_status"`
}

// NewAIRequestResult is a constructor for the ai request result struct
func NewAIRequestResult(
	requestID string,
	results []webSocket.ValResultI,
	status string,
) AIRequestResult {
	return AIRequestResult{
		RequestID: requestID,
		Results:   results,
		Status:    status,
	}
}
