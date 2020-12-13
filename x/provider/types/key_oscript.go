package types

var (
	// OScriptKeyPrefix sets a prefix for oScript key
	OScriptKeyPrefix = "o"
	// OScriptKeyPostfix sets the postfix of the oracle script name
	OScriptKeyPostfix = ".py"
)

// OracleScriptStoreKey returns the key to retrieve a specfic oScript from the store.
func OracleScriptStoreKey(name string) []byte {
	return []byte(OScriptKeyPrefix + name + OScriptKeyPostfix)
}

// OracleScriptStoreKeyString returns the key to retrieve a specfic oScript from the store.
func OracleScriptStoreKeyString(name string) string {
	return OScriptKeyPrefix + name + OScriptKeyPostfix
}
