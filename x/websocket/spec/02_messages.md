<!--
order: 3
-->

# Messages

In this section, we describe all the websocket messages that interact with the states of the module. All state objects specified in this document are located within the [state](./01_state.md#pool) section.

## MsgCreateReport

MsgCreateReport defines message for creating a report by a reporter of a validator

```go
type MsgCreateReport struct {
	RequestID         string                       `json:"request_id"`
	DataSourceResults []exported.DataSourceResultI `json:"data_source_results"`
	TestCaseResults   []exported.TestCaseResultI   `json:"test_case_results"`
	Reporter          Reporter                     `json:"reporter"`
	Fees              sdk.Coins                    `json:"report_fee"`
	AggregatedResult  []byte              		   `json:"aggregated_result"`
	ResultStatus      string                       `json:"result_status"`
}
```

If one of the below conditions occurs, the message will not be accepted by the network:

- The reporter address or name is empty
- The request ID or validator address is empty .
- If there is no data source result / test case result / aggregated result.
- The fee type is invalid.
- The report already exists (by checking whether in that block height, the same request ID is reported twice)
- Unexpected error in adding the report into the store.
- the result status is different from the success and fail status.