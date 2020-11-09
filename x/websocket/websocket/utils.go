package websocket

import provider "github.com/oraichain/orai/x/provider/types"

func getOScriptPath(name string) string {
	return provider.ScriptPath + provider.OracleScriptStoreKeyString(name)
}

func getDSourcePath(name string) string {
	return provider.ScriptPath + provider.DataSourceStoreKeyString(name)
}

func getTCasePath(name string) string {
	return provider.ScriptPath + provider.TestCaseStoreKeyString(name)
}

const (
	// ScriptPath is the path that stores all the script files (oracle scripts, data sources, test cases)
	ScriptPath = "../../../.oraifiles/"
)
