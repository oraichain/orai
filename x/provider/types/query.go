package types

type QueryOracleScriptRequest struct {
	Name string `json:"name"`
}

type QueryDataSourceRequest struct {
	Name string `json:"name"`
}

type QueryTestCaseRequest struct {
	Name string `json:"name"`
}

type QueryOracleScriptsRequest struct {
	Name  string `json:"name"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type QueryDataSourcesRequest struct {
	Name  string `json:"name"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type QueryTestCasesRequest struct {
	Name  string `json:"name"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type MinFeesRequest struct {
	OracleScriptName string `json:"oscript_name"`
	ValNum           int    `json:"val_num"`
}
