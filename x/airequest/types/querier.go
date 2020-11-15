package types

import (
	"fmt"
	"strings"

	provider "github.com/oraichain/orai/x/provider/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query aiDataSources supported by the provider querier. Eg: custom provider query oScript
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryAIRequest    = "aireq"
	QueryAIRequestIDs = "aireqs"
	QueryMinFees      = "min_fees"
)

// QueryResAIRequest resolves a query for an AI request
type QueryResAIRequest struct {
	RequestID        string                   `json:"request_id"`
	Creator          sdk.AccAddress           `json:"request_creator"`
	OracleScriptName string                   `json:"oscript_name"`
	Validators       []sdk.ValAddress         `json:"validator_addrs"`
	BlockHeight      int64                    `json:"block_height"`
	AIDataSources    []provider.AIDataSourceI `json:"data_sources"`
	TestCases        []provider.TestCaseI     `json:"test_cases"`
	Fees             string                   `json:"transaction_fee"`
}

// NewQueryResAIRequest is the constructor for the query ai request
func NewQueryResAIRequest(reqID string, creator sdk.AccAddress, oscriptName string, vals []sdk.ValAddress, blockHeight int64, aiDataSources []provider.AIDataSourceI, testCases []provider.TestCaseI, fees string) QueryResAIRequest {
	return QueryResAIRequest{
		RequestID:        reqID,
		Creator:          creator,
		OracleScriptName: oscriptName,
		Validators:       vals,
		BlockHeight:      blockHeight,
		AIDataSources:    aiDataSources,
		TestCases:        testCases,
		Fees:             fees,
	}
}

// implement fmt.Stringer
func (req QueryResAIRequest) String() string {

	valString := fmt.Sprintln(req.Validators)
	dataSourceString := fmt.Sprintln(req.AIDataSources)
	testCaseString := fmt.Sprintln(req.TestCases)

	return strings.TrimSpace(fmt.Sprintf(`RequestID: %s
	OracleScriptName: %s Validators: %s BlockHeight: %d AIDataSources: %s TestCases: %s Fees: %s`, req.RequestID, req.OracleScriptName, valString, req.BlockHeight, dataSourceString, testCaseString, req.Fees))
}

// QueryResAIRequestIDs Queries all Request IDs
type QueryResAIRequestIDs []string

// implement fmt.Stringer
func (e QueryResAIRequestIDs) String() string {
	return strings.Join(e[:], "\n")
}
