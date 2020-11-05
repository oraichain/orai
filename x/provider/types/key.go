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
	DefaultQueryPage = "1"

	// DefaultQueryLimit sets the default query limit value
	DefaultQueryLimit = "5"

	// ScriptPath is the path that stores all the script files (oracle scripts, data sources, test cases)
	ScriptPath = ".oraifiles/"

	// Denom is the denominator of the currency
	Denom = "orai"
)
