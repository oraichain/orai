<!--
order: 1
-->

# State

## OracleScript

When an individual requests data from the Oraichain's oracle, he calls one of the existing Oracle Scripts on Oraichain. An Oracle Script is a smart contract which will collect data from different AI Data Sources or Oracle Scripts before using test cases given by the user to test those data. If they satisfy the requirements of the test cases, the script will trigger an aggregation phase, and then the aggregated data is stored on the Oraichain for future usage and verification.

```go
type OracleScript struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Description string         `json:"description"`
	MinimumFees sdk.Coins      `json:"minimum_fees"`
	DSources    []string       `json:"data_sources"`
	TCases      []string       `json:"test_cases"`
}
```

Note that DSources and TCases are the list of Data Source and Test Case unique identifiers respectively, and they will be mapped to the actual content of the Oracle Script. When creating / editing an Oracle Script, these identifiers will be read from the Oracle Script smart contract before storing them on-chain.

## AIDataSource

AI Data Source is a basic element in the Oraichain network. It describes a way to provide an AI model for users. An AI Provider can create and register different AI Data Sources in the Oraichain using the message `MsgCreateAIDataSource` sent to the system. The AI Data Source metadata will be stored on-chain.

```go
type AIDataSource struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        sdk.Coins      `json:"transaction_fee"`
	Description string         `json:"description"`
}
```

## TestCase
AI test case is one of the most fundamental elements in the Oraichain system. It consists of a set of inputs and outputs encrypted for privacy purposes. This set is indicated by Oracle Scripts to examine the credibility of the AI Data Sources before actually executing the request.

```go
type TestCase struct {
	Owner       sdk.AccAddress `json:"owner"`
	Name        string         `json:"name"`
	Fees        sdk.Coins      `json:"transaction_fee"`
	Description string         `json:"description"`
}
```

If you notice, the actual source code is smart contract code stored on blockchain, and should keep its size as small as possible. We should only keep the stored data small only to minimize the gas fee and execution time for the nodes.
