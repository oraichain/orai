package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query statements for oracle scripts
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryOracleScript      = "oscript"
	QueryOracleScripts     = "oscripts"
	QueryOracleScriptNames = "onames"
	QueryMinFees           = "min_fees"
)

// QueryResOracleScript resolves a query to a oScript
type QueryResOracleScript struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	MinimumFees sdk.Coins      `json:"minimum_fees"`
	DSources    []string       `json:"data_sources"`
	TCases      []string       `json:"test_cases"`
}

// NewQueryResOracleScript is the constructor for the query oracle script request
func NewQueryResOracleScript(name string, owner sdk.AccAddress, code string, des string, minFees sdk.Coins, ds []string, tc []string) QueryResOracleScript {
	return QueryResOracleScript{
		Name:        name,
		Owner:       owner,
		Code:        code,
		Description: des,
		MinimumFees: minFees,
		DSources:    ds,
		TCases:      tc,
	}
}

func (qrs QueryResOracleScript) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s
Owner: %s
Code: %s Description: %s Minimum Fees: %s Data Sources: %s Test Cases: %s`, qrs.Name, string(qrs.Owner[:]), qrs.Code, qrs.Description, qrs.MinimumFees.String(), qrs.DSources, qrs.TCases))
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
