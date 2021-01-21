<!--
order: 1
-->

# State

## Report

Report is an entity that a validator creates to store information relating to executing the Oracle Script.

```go
type Report struct {
	RequestID         string                       `json:"request_id"` // ID of the request
	DataSourceResults []exported.DataSourceResultI `json:"data_source_results"` // data source results after executing them.
	TestCaseResults   []exported.TestCaseResultI   `json:"test_case_results"` // test case results after executing them.
	BlockHeight       int64                        `json:"block_height"` // block height where the report is stored on Oraichain.
	Fees              sdk.Coins                    `json:"report_fee"` // fees for reporting.
	AggregatedResult  []byte              		   `json:"aggregated_result"` // the aggregated result retrieved from the Oracle Script.
	ResultStatus      string                       `json:"result_status"` // the status of the result, can be either fail or success.
	Reporter          Reporter                     `json:"reporter"` // the account that is used to create the report transaciton.
}
```

## DataSourceResult

DataSourceResult stores the result of an AI Data Source.

```go
type DataSourceResult struct {
	Name   string `json:"data_source"` // name / unique ID of the data source. 
	Result []byte `json:"result"` // result after running the data source.
	Status string `json:"result_status"` // status of the result, which is fail or success.
}
```

If the AI Data Source does not pass the Test Case, the status will be `fail`, otherwise `success`.

## TestCaseResult

TestCaseResult stores information of a test case result

```go
type TestCaseResult struct {
	Name              string                       `json:"test_case"` // test case name or unique ID
	DataSourceResults []exported.DataSourceResultI `json:"data_source_result"` // the list of data source results that have been tested by this test case.
}
```

## Reporter

Reporter is the one who send reports from the validator back to Oraichain created by a validator

```go
type Reporter struct {
	Address   sdk.AccAddress `json:"reporter_address"` // account address of the reporter.
	Name      string         `json:"reporter_name"` // name of the reporter.
	Validator sdk.ValAddress `json:"reporter_validator"` // validator address of this reporter.
}
```

## Validator

Validator mimics the original validator to store information of a validator that executes the Oracle Script

```go
type Validator struct {
	Address     sdk.ValAddress `json:"address"` // the validator address.
	VotingPower int64          `json:"voting_power"` // the total voting power of this validator.
	Status      string         `json:"status"` // the status of the validator (active or inactive)
}
```

## ValResult

ValResult stores the result information from a validator that has executed the oracle script

```go
// ValResult stores the result information from a validator that has executed the oracle script
type ValResult struct {
	Validator    exported.ValidatorI `json:"validator"` // the validator object that mimics the real validator with custom attributes.
	Result       []byte              `json:"result"` // the aggregated result in bytes
	ResultStatus string              `json:"result_status"` // fail or success
}
```

