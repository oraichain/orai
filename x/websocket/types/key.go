package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "websocket"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	// IPFSUrl is he default URL of ipfs gateway to store user images
	IPFSUrl = "http://164.90.180.95:5001/api/v0"

	//IPFSAdd is the path for adding a file onto IPFS
	IPFSAdd = "/add"

	//IPFSCat is the path for retrieving a file from IPFS into the system
	IPFSCat = "/cat"

	// FailedResult represents a data source that cannot pass the test case
	FailedResult = "null"

	// ResultFailure is the fail status of a result after the test case runs
	ResultFailure = "fail"

	// ResultSuccess is the success status of a result after the test case runs
	ResultSuccess = "success"

	// RngSeedSize is the size of the random seed for validator sampling.
	RngSeedSize = 64

	// NumSeedRemoval is the number of bytes an old seed is removed to generate a new one
	NumSeedRemoval = 1
	// Denom is the denominator of the currency
	Denom = "orai"
)

var (

	// ReportKeyPrefix sets a prefix for a report key
	ReportKeyPrefix = "rp"

	// ReporterKeyPrefix sets a prefix for a reporter key
	ReporterKeyPrefix = "rer"

	// StrategyKeyPrefix sets a prefix for a strategy key
	StrategyKeyPrefix = "st"
)

// ReportStoreKey returns the key to retrieve a specfic report from the store.
func ReportStoreKey(requestID string, valAddress string) []byte {
	// buf := append([]byte(ReportKeyPrefix), valAddress...)
	// buf = append(buf, []byte(requestID)...)
	return []byte(ReportKeyPrefix + requestID + valAddress)
}

// ReporterStoreKey returns the key to retrieve a specfic report from the store.
func ReporterStoreKey(valAddress sdk.ValAddress, reporterAddress sdk.AccAddress) []byte {
	// buf := append([]byte(ReporterKeyPrefix), []byte(valAddress)...)
	// buf = append(buf, []byte(reporterAddress)...)
	// return buf
	return []byte(ReporterKeyPrefix + string(valAddress[:]) + string(reporterAddress[:]))
}

// ReportersOfValidatorPrefixKey returns the prefix key to get all reporters of a validator.
func ReportersOfValidatorPrefixKey(val sdk.ValAddress) []byte {
	return append([]byte(ReporterKeyPrefix), val.Bytes()...)
}

// ReportStoreKeyPrefix returns the prefix of report key (used to iterate all the reports of a request)
func ReportStoreKeyPrefix(requestID string) []byte {
	return []byte(ReportKeyPrefix + requestID)
}

// ReportStoreKeyPrefixAll returns the prefix of report key (used to iterate all the reports of all requests)
func ReportStoreKeyPrefixAll() []byte {
	return []byte(ReportKeyPrefix)
}

// ReportStoreKeyString returns the key to retrieve a specfic report from the store.
func ReportStoreKeyString(valAddress []byte, requestID string) string {
	return ReportKeyPrefix + string(valAddress[:]) + requestID
}

// StrategyStoreKey returns the key to retrieve a specfic strategy from the store.
func StrategyStoreKey(stratID uint64, stratName string) []byte {
	return []byte(StrategyKeyPrefix + strconv.FormatUint(stratID, 10) + stratName)
}
