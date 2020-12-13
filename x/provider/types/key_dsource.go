package types

var (

	// DataSourceKeyPrefix sets a prefix for data source key
	DataSourceKeyPrefix = "d"
	// DataSourcePostFix sets the postfix of the data source name
	DataSourcePostFix = ".py"
)

// DataSourceStoreKey returns the key to retrieve a specfic data source from the store.
func DataSourceStoreKey(name string) []byte {
	return []byte(DataSourceKeyPrefix + name + DataSourcePostFix)
}

// DataSourceStoreKeyString returns the key to retrieve a specfic data source from the store.
func DataSourceStoreKeyString(name string) string {
	return DataSourceKeyPrefix + name + DataSourcePostFix
}
