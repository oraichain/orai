<!--
order: 7
-->

# Events

The `provider` module emits some events as follows:

## MsgCreateOracleScript

| Type        | Attribute Key | Attribute Value    |
| ----------- | ------------- | ------------------ |
| set_oscript | oscript_name  | {oracleScriptName} |

## MsgEditOracleScript

| Type         | Attribute Key | Attribute Value    |
| ------------ | ------------- | ------------------ |
| edit_oscript | oscript_name  | {oracleScriptName} |

## MsgCreateAIDataSource

| Type           | Attribute Key   | Attribute Value  |
| -------------- | --------------- | ---------------- |
| set_datasource | datasource_name | {dataSourceName} |

## MsgEditAIDataSource

| Type            | Attribute Key   | Attribute Value  |
| --------------- | --------------- | ---------------- |
| edit_datasource | datasource_name | {dataSourceName} |

## MsgCreateTestCase

| Type          | Attribute Key  | Attribute Value |
| ------------- | -------------- | --------------- |
| set_test_case | test_case_name | {testCaseName}  |

## MsgEditTestCase

| Type           | Attribute Key  | Attribute Value |
| -------------- | -------------- | --------------- |
| edit_test_case | test_case_name | {testCaseName}  |