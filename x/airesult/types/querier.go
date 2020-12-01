package types

import (
	aiRequest "github.com/oraichain/orai/x/airequest/exported"
	webSocket "github.com/oraichain/orai/x/websocket/exported"
)

// Query aiDataSources supported by the provider querier. Eg: custom provider query oScript
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryFullRequest = "fullreq"
	QueryMinFees     = "min_fees"
)

// QueryResFullRequest Queries a complete request with reports
type QueryResFullRequest struct {
	AIRequest aiRequest.AIRequestI `json:"ai_request"`
	Reports   []webSocket.ReportI  `json:"reports"`
	Result    AIRequestResult      `json:"ai_result"`
}

// NewQueryResFullRequest is the constructor for the query full request
func NewQueryResFullRequest(aiReq aiRequest.AIRequestI, reps []webSocket.ReportI, result AIRequestResult) QueryResFullRequest {
	return QueryResFullRequest{
		AIRequest: aiReq,
		Reports:   reps,
		Result:    result,
	}
}
