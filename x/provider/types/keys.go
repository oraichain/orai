package types

const (
	// ModuleName is the name of the module
	ModuleName = "provider"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	// DefaultQueryPage sets the default page query value
	DefaultQueryPage = 1

	// DefaultQueryLimit sets the default query limit value
	DefaultQueryLimit = 5

	// DefaultValNum default validator number
	DefaultValNum = 1

	// ScriptPath is the path that stores all the script files (oracle scripts, data sources, test cases)
	ScriptPath = ".oraifiles/"

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"

	// Denom is the denominator of the currency
	Denom = "orai"

	// DataSourceKeyPrefix sets a prefix for data source key
	DataSourceKeyPrefix = "d"

	// TestCaseKeyPrefix sets a prefix for test case key
	TestCaseKeyPrefix = "tc"

	// OScriptKeyPrefix sets a prefix for oScript key
	OScriptKeyPrefix = "o"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// TestCaseStoreKey returns the key to retrieve a specfic test case from the store.
func TestCaseStoreKey(name string) []byte {
	return []byte(TestCaseKeyPrefix + name)
}

// OracleScriptStoreKey returns the key to retrieve a specfic oScript from the store.
func OracleScriptStoreKey(name string) []byte {
	return []byte(OScriptKeyPrefix + name)
}

// DataSourceStoreKey returns the key to retrieve a specfic data source from the store.
func DataSourceStoreKey(name string) []byte {
	return []byte(DataSourceKeyPrefix + name)
}
