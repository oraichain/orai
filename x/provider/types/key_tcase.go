package types

var (
	// TestCaseKeyPrefix sets a prefix for result key
	TestCaseKeyPrefix = "tc"
)

// TestCaseStoreKey returns the key to retrieve a specfic test case from the store.
func TestCaseStoreKey(name string) []byte {
	return []byte(TestCaseKeyPrefix + name)
}

// TestCaseStoreKeyString returns the key to retrieve a specfic test case from the store.
func TestCaseStoreKeyString(name string) string {
	return TestCaseKeyPrefix + name
}
