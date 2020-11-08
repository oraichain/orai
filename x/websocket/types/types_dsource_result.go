package types

import "github.com/oraichain/orai/x/websocket/exported"

// type check for the implementation of the interface DataSourceResultI
var _ exported.DataSourceResultI = (*DataSourceResult)(nil)

// DataSourceResult stores the data source result
type DataSourceResult struct {
	Name   string `json:"data_source"`
	Result []byte `json:"result"`
	Status string `json:"result_status"`
}

// NewDataSourceResult is the constructor of the data source result struct
func NewDataSourceResult(
	name string,
	result []byte,
	status string,
) DataSourceResult {
	return DataSourceResult{
		Name:   name,
		Result: result,
		Status: status,
	}
}

// GetName is getter method for DataSourceResult struct
func (ds DataSourceResult) GetName() string {
	return ds.Name
}

// GetResult is getter method for DataSourceResult struct
func (ds DataSourceResult) GetResult() []byte {
	return ds.Result
}

// GetStatus is getter method for DataSourceResult struct
func (ds DataSourceResult) GetStatus() string {
	return ds.Status
}
