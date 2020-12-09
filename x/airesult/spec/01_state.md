<!--
order: 1
-->

# State

## AIRequestResult

AIRequestResult stores the final result after aggregating the results from the reports of an AI request

```go
type AIRequestResult struct {
	RequestID string                 `json:"request_id"` // Request ID.
	Results   []webSocket.ValResultI `json:"results"` // Results from different validators.
	Status    string                 `json:"request_status"` // Status of the result (pending / finished / expired)
}
```

## Reward

Reward stores a list of validators, data source owners and test case owners that receive rewards for a specific block height

```go
type Reward struct {
	Validators    []webSocket.ValidatorI   `json:"validators"` // List of validators that have reported successfully.
	DataSources   []provider.AIDataSourceI `json:"data_sources"` // List of data sources that have been executed successfully.
	TestCases     []provider.TestCaseI     `json:"test_cases"` // List of test cases that have been executed successfully.
	BlockHeight   int64                    `json:"block_height"` // The block height that this reward object is stored.
	TotalPower    int64                    `json:"total_voting_power"` // Total voting power of all validators have reported.
	ProviderFees  sdk.Coins                `json:"provider_fees"` // Total reward for the providers.
	ValidatorFees sdk.Coins                `json:"validator_fees"` // Total reward for the validators.
}
```
