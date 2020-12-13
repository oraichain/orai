package types

var (
	// TestCaseKeyPrefix sets a prefix for result key
	TestCaseKeyPrefix = "tc"
	// TestCaseKeyPostfix sets the postfix of the oracle script name
	TestCaseKeyPostfix = ".py"
)

// TestCaseStoreKey returns the key to retrieve a specfic test case from the store.
func TestCaseStoreKey(name string) []byte {
	return []byte(TestCaseKeyPrefix + name + TestCaseKeyPostfix)
}

// TestCaseStoreKeyString returns the key to retrieve a specfic test case from the store.
func TestCaseStoreKeyString(name string) string {
	return TestCaseKeyPrefix + name + TestCaseKeyPostfix
}
