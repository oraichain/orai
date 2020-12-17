<!--
order: 3
-->

# Messages

In this section, we describe all the airequest messages that interact with the states of the module. All state objects specified in this document are located within the [state](./01_state.md#pool) section.

## MsgSetAIRequest

A new AIRequest object is created using the `MsgSetAIRequest` message.

Requests a new data based on an existing Oracle Script. A data request will be assigned a unique identifier with an AI request prefix once the transaction is confirmed. After sufficient validators report successfully. The results of the data requests will be written and stored permanently on Oraichain for future uses.

```go
type MsgSetAIRequest struct {
	RequestID        string         		 `json:"request_id"`
	OracleScriptName string         		 `json:"oscript_name"`
	Creator          sdk.AccAddress 		 `json:"creator"`
	ValidatorCount   int            		 `json:"validator_count"`
	Fees             string         		 `json:"transaction_fee"`
	Input            json.RawMessage         `json:"request_input"`
	ExpectedOutput   json.RawMessage         `json:"expected_output"`
}
```

If one of the below conditions occurs, the message will not be accepted by the network:

- The Oracle Script name is empty.
- The number of validators required is:
  - Smaller or equal to zero.
  - Higher than the maximum number of bonded validators currently running.
  - No validator has voting power.
- The fees are faulty, which are:
  - Invalid data type.
  - The fees are not enough to pay for the validators, Data Source & Test Case owners.
- Cannot get Data Sources & Test Cases from the Oracle Script.
- Cannot get the minimum fees of the Oracle Script.