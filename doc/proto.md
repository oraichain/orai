# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [x/provider/types/genesis.proto](#x/provider/types/genesis.proto)
    - [GenesisState](#oraichain.orai.provider.GenesisState)
  
- [x/provider/types/query.proto](#x/provider/types/query.proto)
    - [Query](#oraichain.orai.provider.Query)
  
- [x/provider/types/query_dsource.proto](#x/provider/types/query_dsource.proto)
    - [DataSourceInfoReq](#oraichain.orai.provider.DataSourceInfoReq)
    - [DataSourceInfoRes](#oraichain.orai.provider.DataSourceInfoRes)
    - [ListDataSourcesReq](#oraichain.orai.provider.ListDataSourcesReq)
    - [ListDataSourcesRes](#oraichain.orai.provider.ListDataSourcesRes)
  
- [x/provider/types/query_oscript.proto](#x/provider/types/query_oscript.proto)
    - [ListOracleScriptsReq](#oraichain.orai.provider.ListOracleScriptsReq)
    - [ListOracleScriptsRes](#oraichain.orai.provider.ListOracleScriptsRes)
    - [MinFeesReq](#oraichain.orai.provider.MinFeesReq)
    - [MinFeesRes](#oraichain.orai.provider.MinFeesRes)
    - [OracleScriptInfoReq](#oraichain.orai.provider.OracleScriptInfoReq)
    - [OracleScriptInfoRes](#oraichain.orai.provider.OracleScriptInfoRes)
  
- [x/provider/types/query_tcase.proto](#x/provider/types/query_tcase.proto)
    - [ListTestCasesReq](#oraichain.orai.provider.ListTestCasesReq)
    - [ListTestCasesRes](#oraichain.orai.provider.ListTestCasesRes)
    - [TestCaseInfoReq](#oraichain.orai.provider.TestCaseInfoReq)
    - [TestCaseInfoRes](#oraichain.orai.provider.TestCaseInfoRes)
  
- [x/provider/types/tx.proto](#x/provider/types/tx.proto)
    - [Msg](#oraichain.orai.provider.Msg)
  
- [x/provider/types/tx_dsource.proto](#x/provider/types/tx_dsource.proto)
    - [MsgCreateAIDataSource](#oraichain.orai.provider.MsgCreateAIDataSource)
    - [MsgCreateAIDataSourceRes](#oraichain.orai.provider.MsgCreateAIDataSourceRes)
    - [MsgEditAIDataSource](#oraichain.orai.provider.MsgEditAIDataSource)
    - [MsgEditAIDataSourceRes](#oraichain.orai.provider.MsgEditAIDataSourceRes)
  
- [x/provider/types/tx_oscript.proto](#x/provider/types/tx_oscript.proto)
    - [MsgCreateOracleScript](#oraichain.orai.provider.MsgCreateOracleScript)
    - [MsgCreateOracleScriptRes](#oraichain.orai.provider.MsgCreateOracleScriptRes)
    - [MsgEditOracleScript](#oraichain.orai.provider.MsgEditOracleScript)
    - [MsgEditOracleScriptRes](#oraichain.orai.provider.MsgEditOracleScriptRes)
  
- [x/provider/types/tx_tcase.proto](#x/provider/types/tx_tcase.proto)
    - [MsgCreateTestCase](#oraichain.orai.provider.MsgCreateTestCase)
    - [MsgCreateTestCaseRes](#oraichain.orai.provider.MsgCreateTestCaseRes)
    - [MsgEditTestCase](#oraichain.orai.provider.MsgEditTestCase)
    - [MsgEditTestCaseRes](#oraichain.orai.provider.MsgEditTestCaseRes)
  
- [x/provider/types/types_ds.proto](#x/provider/types/types_ds.proto)
    - [AIDataSource](#oraichain.orai.provider.AIDataSource)
  
- [x/provider/types/types_os.proto](#x/provider/types/types_os.proto)
    - [OracleScript](#oraichain.orai.provider.OracleScript)
  
- [x/provider/types/types_tc.proto](#x/provider/types/types_tc.proto)
    - [TestCase](#oraichain.orai.provider.TestCase)
  
- [Scalar Value Types](#scalar-value-types)



<a name="x/provider/types/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/genesis.proto



<a name="oraichain.orai.provider.GenesisState"></a>

### GenesisState
GenesisState defines the capability module&#39;s genesis state.


| Field         | Type                                                  | Label    | Description |
| ------------- | ----------------------------------------------------- | -------- | ----------- |
| AIDataSources | [AIDataSource](#oraichain.orai.provider.AIDataSource) | repeated |             |
| OracleScripts | [OracleScript](#oraichain.orai.provider.OracleScript) | repeated |             |
| TestCases     | [TestCase](#oraichain.orai.provider.TestCase)         | repeated |             |





 

 

 

 



<a name="x/provider/types/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/query.proto


 

 

 


<a name="oraichain.orai.provider.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name       | Request Type                                                          | Response Type                                                         | Description                                       |
| ----------------- | --------------------------------------------------------------------- | --------------------------------------------------------------------- | ------------------------------------------------- |
| DataSourceInfo    | [DataSourceInfoReq](#oraichain.orai.provider.DataSourceInfoReq)       | [DataSourceInfoRes](#oraichain.orai.provider.DataSourceInfoRes)       | DataSourceInfo gets the data source meta data     |
| ListDataSources   | [ListDataSourcesReq](#oraichain.orai.provider.ListDataSourcesReq)     | [ListDataSourcesRes](#oraichain.orai.provider.ListDataSourcesRes)     | ListDataSources gets the list of data sources     |
| OracleScriptInfo  | [OracleScriptInfoReq](#oraichain.orai.provider.OracleScriptInfoReq)   | [OracleScriptInfoRes](#oraichain.orai.provider.OracleScriptInfoRes)   | OracleScriptInfo gets the oracle script meta data |
| ListOracleScripts | [ListOracleScriptsReq](#oraichain.orai.provider.ListOracleScriptsReq) | [ListOracleScriptsRes](#oraichain.orai.provider.ListOracleScriptsRes) | ListOracleScripts gets the list of oracle scripts |
| TestCaseInfo      | [TestCaseInfoReq](#oraichain.orai.provider.TestCaseInfoReq)           | [TestCaseInfoRes](#oraichain.orai.provider.TestCaseInfoRes)           | TestCaseInfo gets the test case meta data         |
| ListTestCases     | [ListTestCasesReq](#oraichain.orai.provider.ListTestCasesReq)         | [ListTestCasesRes](#oraichain.orai.provider.ListTestCasesRes)         | ListTestCases gets the list of test cases         |
| QueryMinFees      | [MinFeesReq](#oraichain.orai.provider.MinFeesReq)                     | [MinFeesRes](#oraichain.orai.provider.MinFeesRes)                     | QueryMinFees gets the min fees of oracle script   |

 



<a name="x/provider/types/query_dsource.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/query_dsource.proto



<a name="oraichain.orai.provider.DataSourceInfoReq"></a>

### DataSourceInfoReq
DataSourceInfoReq is the request type for the Query/DataSourceInfo RPC method


| Field | Type              | Label | Description                                     |
| ----- | ----------------- | ----- | ----------------------------------------------- |
| name  | [string](#string) |       | address is the address of the contract to query |






<a name="oraichain.orai.provider.DataSourceInfoRes"></a>

### DataSourceInfoRes
DataSourceInfoRes is the response type for the Query/DataSourceInfo RPC method


| Field       | Type                                                  | Label    | Description                                                                     |
| ----------- | ----------------------------------------------------- | -------- | ------------------------------------------------------------------------------- |
| name        | [string](#string)                                     |          |                                                                                 |
| owner       | [bytes](#bytes)                                       |          | Owner is the address who is allowed to make further changes to the data source. |
| contract    | [string](#string)                                     |          |                                                                                 |
| description | [string](#string)                                     |          |                                                                                 |
| fees        | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |                                                                                 |






<a name="oraichain.orai.provider.ListDataSourcesReq"></a>

### ListDataSourcesReq
ListDataSourcesReq is the request type for the Query/ListDataSources RPC method


| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| name  | [string](#string) |       |             |
| page  | [int64](#int64)   |       |             |
| limit | [int64](#int64)   |       |             |






<a name="oraichain.orai.provider.ListDataSourcesRes"></a>

### ListDataSourcesRes
ListDataSourcesRes is the response type for the Query/ListDataSources RPC method


| Field         | Type                                                  | Label    | Description |
| ------------- | ----------------------------------------------------- | -------- | ----------- |
| AIDataSources | [AIDataSource](#oraichain.orai.provider.AIDataSource) | repeated |             |
| count         | [int64](#int64)                                       |          |             |





 

 

 

 



<a name="x/provider/types/query_oscript.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/query_oscript.proto



<a name="oraichain.orai.provider.ListOracleScriptsReq"></a>

### ListOracleScriptsReq
ListOracleScriptsReq is the request type for the Query/ListOracleScripts RPC method


| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| name  | [string](#string) |       |             |
| page  | [int64](#int64)   |       |             |
| limit | [int64](#int64)   |       |             |






<a name="oraichain.orai.provider.ListOracleScriptsRes"></a>

### ListOracleScriptsRes
ListOracleScriptsRes is the response type for the Query/ListOracleScripts RPC method


| Field         | Type                                                  | Label    | Description |
| ------------- | ----------------------------------------------------- | -------- | ----------- |
| OracleScripts | [OracleScript](#oraichain.orai.provider.OracleScript) | repeated |             |
| count         | [int64](#int64)                                       |          |             |






<a name="oraichain.orai.provider.MinFeesReq"></a>

### MinFeesReq
ListOracleScriptsReq is the request type for the Query/ListOracleScripts RPC method


| Field            | Type              | Label | Description |
| ---------------- | ----------------- | ----- | ----------- |
| OracleScriptName | [string](#string) |       |             |
| ValNum           | [int64](#int64)   |       |             |






<a name="oraichain.orai.provider.MinFeesRes"></a>

### MinFeesRes
ListOracleScriptsRes is the response type for the Query/ListOracleScripts RPC method


| Field   | Type              | Label | Description |
| ------- | ----------------- | ----- | ----------- |
| minFees | [string](#string) |       |             |






<a name="oraichain.orai.provider.OracleScriptInfoReq"></a>

### OracleScriptInfoReq
OracleScriptInfoReq is the request type for the Query/OracleScriptInfo RPC method


| Field | Type              | Label | Description                                     |
| ----- | ----------------- | ----- | ----------------------------------------------- |
| name  | [string](#string) |       | address is the address of the contract to query |






<a name="oraichain.orai.provider.OracleScriptInfoRes"></a>

### OracleScriptInfoRes
OracleScriptInfoRes is the response type for the Query/OracleScriptInfo RPC method


| Field       | Type                                                  | Label    | Description                                                                       |
| ----------- | ----------------------------------------------------- | -------- | --------------------------------------------------------------------------------- |
| name        | [string](#string)                                     |          |                                                                                   |
| owner       | [bytes](#bytes)                                       |          | Owner is the address who is allowed to make further changes to the oracle script. |
| contract    | [string](#string)                                     |          |                                                                                   |
| description | [string](#string)                                     |          |                                                                                   |
| fees        | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |                                                                                   |
| d_sources   | [string](#string)                                     | repeated |                                                                                   |
| t_cases     | [string](#string)                                     | repeated |                                                                                   |





 

 

 

 



<a name="x/provider/types/query_tcase.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/query_tcase.proto



<a name="oraichain.orai.provider.ListTestCasesReq"></a>

### ListTestCasesReq
ListTestCasesReq is the request type for the Query/ListTestCases RPC method


| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| name  | [string](#string) |       |             |
| page  | [int64](#int64)   |       |             |
| limit | [int64](#int64)   |       |             |






<a name="oraichain.orai.provider.ListTestCasesRes"></a>

### ListTestCasesRes
ListTestCasesRes is the response type for the Query/ListTestCases RPC method


| Field     | Type                                          | Label    | Description |
| --------- | --------------------------------------------- | -------- | ----------- |
| TestCases | [TestCase](#oraichain.orai.provider.TestCase) | repeated |             |
| count     | [int64](#int64)                               |          |             |






<a name="oraichain.orai.provider.TestCaseInfoReq"></a>

### TestCaseInfoReq
TestCaseInfoReq is the request type for the Query/TestCaseInfo RPC method


| Field | Type              | Label | Description                                     |
| ----- | ----------------- | ----- | ----------------------------------------------- |
| name  | [string](#string) |       | address is the address of the contract to query |






<a name="oraichain.orai.provider.TestCaseInfoRes"></a>

### TestCaseInfoRes
TestCaseInfoRes is the response type for the Query/TestCaseInfo RPC method


| Field       | Type                                                  | Label    | Description                                                                   |
| ----------- | ----------------------------------------------------- | -------- | ----------------------------------------------------------------------------- |
| name        | [string](#string)                                     |          |                                                                               |
| owner       | [bytes](#bytes)                                       |          | Owner is the address who is allowed to make further changes to the test case. |
| contract    | [string](#string)                                     |          |                                                                               |
| description | [string](#string)                                     |          |                                                                               |
| fees        | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |                                                                               |





 

 

 

 



<a name="x/provider/types/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/tx.proto


 

 

 


<a name="oraichain.orai.provider.Msg"></a>

### Msg
Msg defines the provider Msg service.

| Method Name        | Request Type                                                            | Response Type                                                                 | Description                    |
| ------------------ | ----------------------------------------------------------------------- | ----------------------------------------------------------------------------- | ------------------------------ |
| CreateAIDataSource | [MsgCreateAIDataSource](#oraichain.orai.provider.MsgCreateAIDataSource) | [MsgCreateAIDataSourceRes](#oraichain.orai.provider.MsgCreateAIDataSourceRes) | Create a new data source       |
| EditAIDataSource   | [MsgEditAIDataSource](#oraichain.orai.provider.MsgEditAIDataSource)     | [MsgEditAIDataSourceRes](#oraichain.orai.provider.MsgEditAIDataSourceRes)     | Edit an existing data source   |
| CreateOracleScript | [MsgCreateOracleScript](#oraichain.orai.provider.MsgCreateOracleScript) | [MsgCreateOracleScriptRes](#oraichain.orai.provider.MsgCreateOracleScriptRes) | Create a new oracle script     |
| EditOracleScript   | [MsgEditOracleScript](#oraichain.orai.provider.MsgEditOracleScript)     | [MsgEditOracleScriptRes](#oraichain.orai.provider.MsgEditOracleScriptRes)     | Edit an existing oracle script |
| CreateTestCase     | [MsgCreateTestCase](#oraichain.orai.provider.MsgCreateTestCase)         | [MsgCreateTestCaseRes](#oraichain.orai.provider.MsgCreateTestCaseRes)         | Create a new test case         |
| EditTestCase       | [MsgEditTestCase](#oraichain.orai.provider.MsgEditTestCase)             | [MsgEditTestCaseRes](#oraichain.orai.provider.MsgEditTestCaseRes)             | Edit an existing test case     |

 



<a name="x/provider/types/tx_dsource.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/tx_dsource.proto



<a name="oraichain.orai.provider.MsgCreateAIDataSource"></a>

### MsgCreateAIDataSource
MsgCreateAIDataSource submit data source metadata onto Oraichain


| Field       | Type              | Label | Description                                                                     |
| ----------- | ----------------- | ----- | ------------------------------------------------------------------------------- |
| name        | [string](#string) |       |                                                                                 |
| description | [string](#string) |       |                                                                                 |
| contract    | [string](#string) |       |                                                                                 |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the data source. |
| fees        | [string](#string) |       |                                                                                 |






<a name="oraichain.orai.provider.MsgCreateAIDataSourceRes"></a>

### MsgCreateAIDataSourceRes
MsgCreateAIDataSourceRes returns store result data.


| Field       | Type              | Label | Description                                                                     |
| ----------- | ----------------- | ----- | ------------------------------------------------------------------------------- |
| name        | [string](#string) |       |                                                                                 |
| description | [string](#string) |       |                                                                                 |
| contract    | [string](#string) |       |                                                                                 |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the data source. |
| fees        | [string](#string) |       |                                                                                 |






<a name="oraichain.orai.provider.MsgEditAIDataSource"></a>

### MsgEditAIDataSource
MsgEditAIDataSource edit data source metadata onto Oraichain


| Field       | Type              | Label | Description                                                                     |
| ----------- | ----------------- | ----- | ------------------------------------------------------------------------------- |
| old_name    | [string](#string) |       |                                                                                 |
| new_name    | [string](#string) |       |                                                                                 |
| description | [string](#string) |       |                                                                                 |
| contract    | [string](#string) |       |                                                                                 |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the data source. |
| fees        | [string](#string) |       |                                                                                 |






<a name="oraichain.orai.provider.MsgEditAIDataSourceRes"></a>

### MsgEditAIDataSourceRes
MsgEditAIDataSourceRes returns edit result data.


| Field       | Type              | Label | Description                                                                     |
| ----------- | ----------------- | ----- | ------------------------------------------------------------------------------- |
| name        | [string](#string) |       |                                                                                 |
| description | [string](#string) |       |                                                                                 |
| contract    | [string](#string) |       |                                                                                 |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the data source. |
| fees        | [string](#string) |       |                                                                                 |





 

 

 

 



<a name="x/provider/types/tx_oscript.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/tx_oscript.proto



<a name="oraichain.orai.provider.MsgCreateOracleScript"></a>

### MsgCreateOracleScript
MsgCreateOracleScript submit oracle script metadata onto Oraichain


| Field        | Type              | Label    | Description                                                                       |
| ------------ | ----------------- | -------- | --------------------------------------------------------------------------------- |
| name         | [string](#string) |          |                                                                                   |
| description  | [string](#string) |          |                                                                                   |
| contract     | [string](#string) |          |                                                                                   |
| owner        | [bytes](#bytes)   |          | Owner is the address who is allowed to make further changes to the oracle script. |
| fees         | [string](#string) |          |                                                                                   |
| data_sources | [string](#string) | repeated |                                                                                   |
| test_cases   | [string](#string) | repeated |                                                                                   |






<a name="oraichain.orai.provider.MsgCreateOracleScriptRes"></a>

### MsgCreateOracleScriptRes
MsgCreateOracleScriptRes returns store result data.


| Field        | Type              | Label    | Description                                                                       |
| ------------ | ----------------- | -------- | --------------------------------------------------------------------------------- |
| name         | [string](#string) |          |                                                                                   |
| description  | [string](#string) |          |                                                                                   |
| contract     | [string](#string) |          |                                                                                   |
| owner        | [bytes](#bytes)   |          | Owner is the address who is allowed to make further changes to the oracle script. |
| fees         | [string](#string) |          |                                                                                   |
| data_sources | [string](#string) | repeated |                                                                                   |
| test_cases   | [string](#string) | repeated |                                                                                   |






<a name="oraichain.orai.provider.MsgEditOracleScript"></a>

### MsgEditOracleScript
MsgEditOracleScript edit oracle script metadata onto Oraichain


| Field        | Type              | Label    | Description                                                                       |
| ------------ | ----------------- | -------- | --------------------------------------------------------------------------------- |
| old_name     | [string](#string) |          |                                                                                   |
| new_name     | [string](#string) |          |                                                                                   |
| description  | [string](#string) |          |                                                                                   |
| contract     | [string](#string) |          |                                                                                   |
| owner        | [bytes](#bytes)   |          | Owner is the address who is allowed to make further changes to the oracle script. |
| fees         | [string](#string) |          |                                                                                   |
| data_sources | [string](#string) | repeated |                                                                                   |
| test_cases   | [string](#string) | repeated |                                                                                   |






<a name="oraichain.orai.provider.MsgEditOracleScriptRes"></a>

### MsgEditOracleScriptRes
MsgEditOracleScriptRes returns edit result data.


| Field        | Type              | Label    | Description                                                                       |
| ------------ | ----------------- | -------- | --------------------------------------------------------------------------------- |
| name         | [string](#string) |          |                                                                                   |
| description  | [string](#string) |          |                                                                                   |
| contract     | [string](#string) |          |                                                                                   |
| owner        | [bytes](#bytes)   |          | Owner is the address who is allowed to make further changes to the oracle script. |
| fees         | [string](#string) |          |                                                                                   |
| data_sources | [string](#string) | repeated |                                                                                   |
| test_cases   | [string](#string) | repeated |                                                                                   |





 

 

 

 



<a name="x/provider/types/tx_tcase.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/tx_tcase.proto



<a name="oraichain.orai.provider.MsgCreateTestCase"></a>

### MsgCreateTestCase
MsgCreateTestCase submit test case metadata onto Oraichain


| Field       | Type              | Label | Description                                                                   |
| ----------- | ----------------- | ----- | ----------------------------------------------------------------------------- |
| name        | [string](#string) |       |                                                                               |
| description | [string](#string) |       |                                                                               |
| contract    | [string](#string) |       |                                                                               |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the test case. |
| fees        | [string](#string) |       |                                                                               |






<a name="oraichain.orai.provider.MsgCreateTestCaseRes"></a>

### MsgCreateTestCaseRes
MsgCreateTestCaseRes returns store result data.


| Field       | Type              | Label | Description                                                                   |
| ----------- | ----------------- | ----- | ----------------------------------------------------------------------------- |
| name        | [string](#string) |       |                                                                               |
| description | [string](#string) |       |                                                                               |
| contract    | [string](#string) |       |                                                                               |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the test case. |
| fees        | [string](#string) |       |                                                                               |






<a name="oraichain.orai.provider.MsgEditTestCase"></a>

### MsgEditTestCase
MsgEditTestCase edit test case metadata onto Oraichain


| Field       | Type              | Label | Description                                                                   |
| ----------- | ----------------- | ----- | ----------------------------------------------------------------------------- |
| old_name    | [string](#string) |       |                                                                               |
| new_name    | [string](#string) |       |                                                                               |
| description | [string](#string) |       |                                                                               |
| contract    | [string](#string) |       |                                                                               |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the test case. |
| fees        | [string](#string) |       |                                                                               |






<a name="oraichain.orai.provider.MsgEditTestCaseRes"></a>

### MsgEditTestCaseRes
MsgEditTestCaseRes returns edit result data.


| Field       | Type              | Label | Description                                                                   |
| ----------- | ----------------- | ----- | ----------------------------------------------------------------------------- |
| name        | [string](#string) |       |                                                                               |
| description | [string](#string) |       |                                                                               |
| contract    | [string](#string) |       |                                                                               |
| owner       | [bytes](#bytes)   |       | Owner is the address who is allowed to make further changes to the test case. |
| fees        | [string](#string) |       |                                                                               |





 

 

 

 



<a name="x/provider/types/types_ds.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/types_ds.proto



<a name="oraichain.orai.provider.AIDataSource"></a>

### AIDataSource



| Field       | Type                                                  | Label    | Description                                                                     |
| ----------- | ----------------------------------------------------- | -------- | ------------------------------------------------------------------------------- |
| name        | [string](#string)                                     |          |                                                                                 |
| contract    | [string](#string)                                     |          |                                                                                 |
| owner       | [bytes](#bytes)                                       |          | Owner is the address who is allowed to make further changes to the data source. |
| description | [string](#string)                                     |          |                                                                                 |
| fees        | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |                                                                                 |





 

 

 

 



<a name="x/provider/types/types_os.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/types_os.proto



<a name="oraichain.orai.provider.OracleScript"></a>

### OracleScript



| Field        | Type                                                  | Label    | Description                                                                     |
| ------------ | ----------------------------------------------------- | -------- | ------------------------------------------------------------------------------- |
| name         | [string](#string)                                     |          |                                                                                 |
| contract     | [string](#string)                                     |          |                                                                                 |
| owner        | [bytes](#bytes)                                       |          | Owner is the address who is allowed to make further changes to the data source. |
| description  | [string](#string)                                     |          |                                                                                 |
| minimum_fees | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |                                                                                 |
| d_sources    | [string](#string)                                     | repeated |                                                                                 |
| t_cases      | [string](#string)                                     | repeated |                                                                                 |





 

 

 

 



<a name="x/provider/types/types_tc.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/provider/types/types_tc.proto



<a name="oraichain.orai.provider.TestCase"></a>

### TestCase



| Field       | Type                                                  | Label    | Description                                                                     |
| ----------- | ----------------------------------------------------- | -------- | ------------------------------------------------------------------------------- |
| name        | [string](#string)                                     |          |                                                                                 |
| contract    | [string](#string)                                     |          |                                                                                 |
| owner       | [bytes](#bytes)                                       |          | Owner is the address who is allowed to make further changes to the data source. |
| description | [string](#string)                                     |          |                                                                                 |
| fees        | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |                                                                                 |





 

 

 

 



## Scalar Value Types

| .proto Type                    | Notes                                                                                                                                           | C++    | Java       | Python      | Go      | C#         | PHP            | Ruby                           |
| ------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ---------- | ----------- | ------- | ---------- | -------------- | ------------------------------ |
| <a name="double" /> double     |                                                                                                                                                 | double | double     | float       | float64 | double     | float          | Float                          |
| <a name="float" /> float       |                                                                                                                                                 | float  | float      | float       | float32 | float      | float          | Float                          |
| <a name="int32" /> int32       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="int64" /> int64       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="uint32" /> uint32     | Uses variable-length encoding.                                                                                                                  | uint32 | int        | int/long    | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64     | Uses variable-length encoding.                                                                                                                  | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32     | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s.                            | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64     | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s.                            | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="fixed32" /> fixed32   | Always four bytes. More efficient than uint32 if values are often greater than 2^28.                                                            | uint32 | int        | int         | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64   | Always eight bytes. More efficient than uint64 if values are often greater than 2^56.                                                           | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum                         |
| <a name="sfixed32" /> sfixed32 | Always four bytes.                                                                                                                              | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes.                                                                                                                             | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="bool" /> bool         |                                                                                                                                                 | bool   | boolean    | boolean     | bool    | bool       | boolean        | TrueClass/FalseClass           |
| <a name="string" /> string     | A string must always contain UTF-8 encoded or 7-bit ASCII text.                                                                                 | string | String     | str/unicode | string  | string     | string         | String (UTF-8)                 |
| <a name="bytes" /> bytes       | May contain any arbitrary sequence of bytes.                                                                                                    | string | ByteString | str         | []byte  | ByteString | string         | String (ASCII-8BIT)            |

