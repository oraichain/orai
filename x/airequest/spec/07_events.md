<!--
order: 7
-->

# Events

The `airequest` module emits some events as follows:

## MsgSetAIRequest

| Type            | Attribute Key           | Attribute Value         |
| --------------- | ----------------------- | ----------------------- |
| set_ai_request  | request_id              | {requestID}             |
| set_ai_request  | request_validator       | {validatorAddress}      |
| ai_request_data | request_id              | {requestID}             |
| ai_request_data | oscript_name            | {oScriptName}           |
| ai_request_data | request_creator         | {accAddress}            |
| ai_request_data | request_input           | {msgInput}              |
| ai_request_data | expected_request_output | {msgExpectedOuput}      |
| ai_request_data | data_sources            | {dataSourceIdentifiers} |
| ai_request_data | test_cases              | {testCaseIdentifiers}   |
