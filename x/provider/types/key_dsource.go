package types

var (

	// DataSourceKeyPrefix sets a prefix for data source key
	DataSourceKeyPrefix = "d"
	// DataSourceFileSuffix sets the suffix of a data source file to python file
	DataSourceFileSuffix = ".py"
)

// DataSourceStoreKey returns the key to retrieve a specfic data source from the store.
func DataSourceStoreKey(name string) []byte {
	return []byte(DataSourceKeyPrefix + name)
}

// DataSourceStoreKeyString returns the key to retrieve a specfic data source from the store.
func DataSourceStoreKeyString(name string) string {
	return DataSourceKeyPrefix + name
}

// DataSourceStoreFileString returns the full key with .py to retrieve a specfic data source from the store.
func DataSourceStoreFileString(name string) string {
	return DataSourceKeyPrefix + name + DataSourceFileSuffix
}
