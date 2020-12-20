package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query statements for data sources
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryDataSource      = "datasource"
	QueryDataSources     = "datasources"
	QueryDataSourceNames = "dnames"
)

// QueryResAIDataSource resolves a query to an data source
type QueryResAIDataSource struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	Fees        sdk.Coins      `json:"fees"`
}

// NewQueryResAIDataSource is the constructor for the query ai data source request
func NewQueryResAIDataSource(name string, owner sdk.AccAddress, code string, des string, fees sdk.Coins) QueryResAIDataSource {
	return QueryResAIDataSource{
		Name:        name,
		Owner:       owner,
		Code:        code,
		Description: des,
		Fees:        fees,
	}
}

func (qre QueryResAIDataSource) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s
Owner: %s
Code: %s Description: %s fees: %s`, qre.Name, string(qre.Owner[:]), qre.Code, qre.Description, qre.Fees))
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

// QueryResAIDataSourceNames Queries Result Payload for a names query
type QueryResAIDataSourceNames []string

// implement fmt.Stringer
func (e QueryResAIDataSourceNames) String() string {
	return strings.Join(e[:], "\n")
}
