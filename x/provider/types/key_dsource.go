package types

var (

	// DataSourceKeyPrefix sets a prefix for data source key
	DataSourceKeyPrefix = "d"
)

// DataSourceStoreKey returns the key to retrieve a specfic data source from the store.
func DataSourceStoreKey(name string) []byte {
	return []byte(DataSourceKeyPrefix + name)
}
