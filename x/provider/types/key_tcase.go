package types

var (
	// TestCaseKeyPrefix sets a prefix for test case key
	TestCaseKeyPrefix = "tc"
	// TestCaseKeySuffix sets a suffix for test case file being python
	TestCaseFileSuffix = ".py"
)

// TestCaseStoreKey returns the key to retrieve a specfic test case from the store.
func TestCaseStoreKey(name string) []byte {
	return []byte(TestCaseKeyPrefix + name)
}

// TestCaseStoreKeyString returns the key to retrieve a specfic test case from the store.
func TestCaseStoreKeyString(name string) string {
	return TestCaseKeyPrefix + name
}

// TestCaseStoreFileString returns the complete file name to retrieve a specfic test case from the store.
func TestCaseStoreFileString(name string) string {
	return TestCaseKeyPrefix + name + TestCaseFileSuffix
}
