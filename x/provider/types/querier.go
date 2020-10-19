package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query aiDataSources supported by the provider querier. Eg: custom provider query oScript
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryOracleScript      = "oscript"
	QueryDataSource        = "datasource"
	QueryDataSources       = "datasources"
	QueryOracleScriptNames = "onames"
	QueryOracleScripts     = "oscripts"
	QueryDataSourceNames   = "dnames"
	QueryAIRequest         = "aireq"
	QueryTestCase          = "testcase"
	QueryTestCases         = "testcases"
	QueryAIRequestIDs      = "aireqs"
	QueryTestCaseNames     = "tcnames"
	QueryFullRequest       = "fullreq"
	QueryMinFees           = "min_fees"
)

// QueryResOracleScript resolves a query to a oScript
type QueryResOracleScript struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	MinimumFees sdk.Coins      `json:"minimum_fees"`
}

// NewQueryResOracleScript is the constructor for the query oracle script request
func NewQueryResOracleScript(name string, owner sdk.AccAddress, code string, des string, minFees sdk.Coins) QueryResOracleScript {
	return QueryResOracleScript{
		Name:        name,
		Owner:       owner,
		Code:        code,
		Description: des,
		MinimumFees: minFees,
	}
}

func (qrs QueryResOracleScript) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s
Owner: %s
Code: %s Description: %s Minimum Fees: %s`, qrs.Name, string(qrs.Owner[:]), qrs.Code, qrs.Description, qrs.MinimumFees.String()))
}

// QueryResAIDataSource resolves a query to an data source
type QueryResAIDataSource struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
}

// NewQueryResAIDataSource is the constructor for the query ai data source request
func NewQueryResAIDataSource(name string, owner sdk.AccAddress, code string, des string) QueryResAIDataSource {
	return QueryResAIDataSource{
		Name:        name,
		Owner:       owner,
		Code:        code,
		Description: des,
	}
}

func (qre QueryResAIDataSource) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s
Owner: %s
Code: %s Description: %s`, qre.Name, string(qre.Owner[:]), qre.Code, qre.Description))
}

// QueryResAIDataSources Queries the list of data sources
type QueryResAIDataSources struct {
	QueryResAIDSources []QueryResAIDataSource `json:"data_sources"`
	Count              int                    `json:"count"`
}

// NewQueryResAIDataSources is the constructor for the query ai data source request
func NewQueryResAIDataSources(queryResAIDataSources []QueryResAIDataSource, count int) QueryResAIDataSources {
	return QueryResAIDataSources{
		QueryResAIDSources: queryResAIDataSources,
		Count:              count,
	}
}

// QueryResOracleScriptNames Queries Result Payload for a names query
type QueryResOracleScriptNames []string

// QueryResOracleScripts Queries the list of oracle scripts
type QueryResOracleScripts struct {
	QueryResOScripts []QueryResOracleScript `json:"oracle_scripts"`
	Count            int                    `json:"count"`
}

// NewQueryResOracleScripts is the constructor for the query oracle scripts request
func NewQueryResOracleScripts(queryResOracleScripts []QueryResOracleScript, count int) QueryResOracleScripts {
	return QueryResOracleScripts{
		QueryResOScripts: queryResOracleScripts,
		Count:            count,
	}
}

// implement fmt.Stringer
func (c QueryResOracleScriptNames) String() string {
	return strings.Join(c[:], "\n")
}

// QueryResAIDataSourceNames Queries Result Payload for a names query
type QueryResAIDataSourceNames []string

// implement fmt.Stringer
func (e QueryResAIDataSourceNames) String() string {
	return strings.Join(e[:], "\n")
}

// QueryResAIRequest resolves a query for an AI request
type QueryResAIRequest struct {
	RequestID        string           `json:"request_id"`
	Creator          sdk.AccAddress   `json:"request_creator"`
	OracleScriptName string           `json:"oscript_name"`
	Validators       []sdk.ValAddress `json:"validator_addrs"`
	BlockHeight      int64            `json:"block_height"`
	AIDataSources    []AIDataSource   `json:"data_sources"`
	TestCases        []TestCase       `json:"test_cases"`
	Fees             string           `json:"transaction_fee"`
}

// NewQueryResAIRequest is the constructor for the query ai request
func NewQueryResAIRequest(reqID string, creator sdk.AccAddress, oscriptName string, vals []sdk.ValAddress, blockHeight int64, aiDataSources []AIDataSource, testCases []TestCase, fees string) QueryResAIRequest {
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

// QueryResTestCase resolves a query for an AI request result
type QueryResTestCase struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
}

// NewQueryResTestCase is the constructor for the query test case request
func NewQueryResTestCase(name string, owner sdk.AccAddress, code string, des string) QueryResTestCase {
	return QueryResTestCase{
		Name:        name,
		Owner:       owner,
		Code:        code,
		Description: des,
	}
}

// implement fmt.Stringer
func (tc QueryResTestCase) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s Owner: %s Code: %s Description: %s`, tc.Name, string(tc.Owner[:]), tc.Code, tc.Description))
}

// QueryResTestCases Queries the list of test cases
type QueryResTestCases struct {
	QueryResTCases []QueryResTestCase `json:"test_cases"`
	Count          int                `json:"count"`
}

// NewQueryResTestCases is the constructor for the query test cases
func NewQueryResTestCases(queryResTestCases []QueryResTestCase, count int) QueryResTestCases {
	return QueryResTestCases{
		QueryResTCases: queryResTestCases,
		Count:          count,
	}
}

// QueryResTestCaseNames Queries all test case IDs
type QueryResTestCaseNames []string

// implement fmt.Stringer
func (req QueryResTestCaseNames) String() string {
	return strings.Join(req[:], "\n")
}

// QueryResFullRequest Queries a complete request with reports
type QueryResFullRequest struct {
	AIRequest AIRequest       `json:"ai_request"`
	Reports   []Report        `json:"reports"`
	Result    AIRequestResult `json:"ai_result"`
}

// NewQueryResFullRequest is the constructor for the query full request
func NewQueryResFullRequest(aiReq AIRequest, reps []Report, result AIRequestResult) QueryResFullRequest {
	return QueryResFullRequest{
		AIRequest: aiReq,
		Reports:   reps,
		Result:    result,
	}
}

// QueryResMinFees Queries a minimum fee value for an oracle script
type QueryResMinFees struct {
	MinFees string `json:"minimum_fees"`
}

// NewQueryResMinFees is the constructor for the query minimum fees
func NewQueryResMinFees(minFees string) QueryResMinFees {
	return QueryResMinFees{
		MinFees: minFees,
	}
}

/*
Below you will be able how to set your own queries:


// QueryResList Queries Result Payload for a query
type QueryResList []string

// implement fmt.Stringer
func (n QueryResList) String() string {
	return strings.Join(n[:], "\n")
}

*/
