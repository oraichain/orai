<!--
order: 1
-->

# State

## AIRequest

AIRequest is the main entry for token holders who want to interact with the AI services provided by AI developers all around the world.

```go
type AIRequest struct {
	RequestID        string                   `json:"request_id"` // The unique identifier of this oracle request
	OracleScriptName string                   `json:"oscript_name"` // The unique name identifier of the Oracle Script
	Creator          sdk.AccAddress           `json:"request_creator"` // The address of the message's sender or creator of this request
	Validators       []sdk.ValAddress         `json:"validator_addr"` // The addresses of the validators that participate in the request
	BlockHeight      int64                    `json:"block_height"` // request block height
	AIDataSources    []provider.AIDataSourceI `json:"data_sources"` // List of AI Data Sources that are in the Oracle Script
	TestCases        []provider.TestCaseI     `json:"testcases"` // List of Test Cases that are in the Oracle Script
	Fees             sdk.Coins                `json:"transaction_fee"` // The transaction fee required to run this AI Request. Eg: 5000orai
	Input            []byte                   `json:"request_input"` // User's input for the AI Request
	ExpectedOutput   []byte                   `json:"expected_output"` // User's expected output for the AI Request
}
```

Input is reserved for passing data to the AI Data Source, while the ExpectedOuput is for the Test Case to use. They are both byte arrays, which make them generic enough for arbitrary objects or arrays to be encoded and passed through. 