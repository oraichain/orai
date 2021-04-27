package types

import "strconv"

const (
	// ModuleName is the name of the module
	ModuleName = "airequest"

	// Denom is the coin denom used for the module
	Denom = "orai"

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

	// RngSeedSize is the size of the random seed for validator sampling.
	RngSeedSize = 64

	// NumSeedRemoval is the number of bytes an old seed is removed to generate a new one
	NumSeedRemoval = 1

	// FailedResponseOs represents an oracle script cannot collect or aggregate the data source results
	FailedResponseOs = "The oracle script could not collect or aggregate the data source results"

	// FailedResponseTc represents the failed response from the data source
	FailedResponseTc = "This data source did not pass the test case given the user expected output"

	// FailedResponseDs represents the failed response from the data source
	FailedResponseDs = "This data source passed the test case but failed when actually running"

	// ResultFailure is the fail status of a result after the test case runs
	ResultFailure = "fail"

	// ResultSuccess is the success status of a result after the test case runs
	ResultSuccess = "success"

	// RequestStatusPending is the status pending of the request
	RequestStatusPending = "pending"

	// RequestStatusFinished is the status finished of the request
	RequestStatusFinished = "finished"

	// RequestStatusExpired is the status expired of the request
	RequestStatusExpired = "expired"

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (

	// RequestKeyPrefix sets a prefix for request key
	RequestKeyPrefix = "req"

	// SeedKeyPrefix sets a prefix for the rng seed
	SeedKeyPrefix = "sd"

	// ReportKeyPrefix sets a prefix for a report key
	ReportKeyPrefix = "rp"

	// TestCaseReportKeyPrefix sets a prefix for a report key
	TestCaseReportKeyPrefix = "tcrp"

	// StrategyKeyPrefix sets a prefix for a strategy key
	StrategyKeyPrefix = "st"

	// ResultKeyPrefix sets a prefix for a result key
	ResultKeyPrefix = "res"

	// RewardKeyPrefix sets a prefix for a reward key
	RewardKeyPrefix = "rw"

	// TestCaseRewardKeyPrefix sets a prefix for a test case reward key
	TestCaseRewardKeyPrefix = "tcrw"
)

// RequestStoreKey returns the key to retrieve a specfic request from the store.
func RequestStoreKey(id string) []byte {
	return []byte(RequestKeyPrefix + id)
}

// RequestStoreKeyString returns the key to retrieve a specfic request from the store.
func RequestStoreKeyString(id string) string {
	return RequestKeyPrefix + id
}

// SeedStoreKey returns the key to retrieve a specfic seed from the store.
func SeedStoreKey() []byte {
	return []byte(SeedKeyPrefix + "seed")
}

// RequeststoreKeyPrefixAll returns the prefix of request key (used to iterate all the requests of all requests
func RequeststoreKeyPrefixAll() []byte {
	return []byte(RequestKeyPrefix)
}

// ReportStoreKey returns the key to retrieve a specfic report from the store.
func ReportStoreKey(requestID string, valAddress string) []byte {
	// buf := append([]byte(ReportKeyPrefix), valAddress...)
	// buf = append(buf, []byte(requestID)...)
	return []byte(ReportKeyPrefix + requestID + valAddress)
}

// ReportStoreKeyPrefix returns the prefix of report key (used to iterate all the reports of a request)
func ReportStoreKeyPrefix(requestID string) []byte {
	return []byte(ReportKeyPrefix + requestID)
}

// ReportStoreKeyPrefixAll returns the prefix of report key (used to iterate all the reports of all requests)
func ReportStoreKeyPrefixAll() []byte {
	return []byte(ReportKeyPrefix)
}

// TestCaseReportStoreKey returns the key to retrieve a specfic report from the store.
func TestCaseReportStoreKey(requestID string, valAddress string) []byte {
	// buf := append([]byte(ReportKeyPrefix), valAddress...)
	// buf = append(buf, []byte(requestID)...)
	return []byte(TestCaseReportKeyPrefix + requestID + valAddress)
}

// TestCaseReportStoreKeyPrefix returns the prefix of report key (used to iterate all the reports of a request)
func TestCaseReportStoreKeyPrefix(requestID string) []byte {
	return []byte(TestCaseReportKeyPrefix + requestID)
}

// ReportStoreKeyPrefixAll returns the prefix of report key (used to iterate all the reports of all requests)
func TestCaseReportStoreKeyPrefixAll() []byte {
	return []byte(TestCaseReportKeyPrefix)
}

// ResultStoreKey returns the key to retrieve a specfic result from the store.
func ResultStoreKey(id string) []byte {
	return []byte(ResultKeyPrefix + id)
}

// RewardStoreKey returns the key to retrieve a specfic reward from the store.
func RewardStoreKey(blockHeight int64) []byte {
	//return []byte(RewardKeyPrefix + string(blockHeight))
	return []byte(RewardKeyPrefix + strconv.FormatInt(blockHeight, 10))
}
