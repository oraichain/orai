package types

var (
	// OScriptKeyPrefix sets a prefix for oScript key
	OScriptKeyPrefix = "o"
)

// OracleScriptStoreKey returns the key to retrieve a specfic oScript from the store.
func OracleScriptStoreKey(name string) []byte {
	return []byte(OScriptKeyPrefix + name)
}

// OracleScriptStoreKeyString returns the key to retrieve a specfic oScript from the store.
func OracleScriptStoreKeyString(name string) string {
	return OScriptKeyPrefix + name
}
