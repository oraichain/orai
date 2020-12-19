package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query statements for test cases
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryTestCase      = "testcase"
	QueryTestCases     = "testcases"
	QueryTestCaseNames = "tcnames"
)

// QueryResTestCase resolves a query for an AI request result
type QueryResTestCase struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	Fees        sdk.Coins      `json:"fees"`
}

// NewQueryResTestCase is the constructor for the query test case request
func NewQueryResTestCase(name string, owner sdk.AccAddress, code string, des string, fees sdk.Coins) QueryResTestCase {
	return QueryResTestCase{
		Name:        name,
		Owner:       owner,
		Code:        code,
		Description: des,
		Fees:        fees,
	}
}

// implement fmt.Stringer
func (tc QueryResTestCase) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s Owner: %s Code: %s Description: %s Fees: %s`, tc.Name, string(tc.Owner[:]), tc.Code, tc.Description, tc.Fees))
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
