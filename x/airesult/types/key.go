package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "airesult"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	// DefaultQueryPage sets the default page query value
	DefaultQueryPage = "1"

	// DefaultQueryLimit sets the default query limit value
	DefaultQueryLimit = "5"

	// RequestStatusPending is the status pending of the request
	RequestStatusPending = "pending"

	// RequestStatusFinished is the status finished of the request
	RequestStatusFinished = "finished"

	// RequestStatusExpired is the status expired of the request
	RequestStatusExpired = "expired"

	// Denom is the denominator of the currency
	Denom = "orai"
)

var (

	// RequestKeyPrefix sets a prefix for request key
	RequestKeyPrefix = "req"

	// ResultKeyPrefix sets a prefix for a result key
	ResultKeyPrefix = "res"

	// ReportKeyPrefix sets a prefix for a report key
	ReportKeyPrefix = "rp"

	// ReporterKeyPrefix sets a prefix for a reporter key
	ReporterKeyPrefix = "rer"

	// RewardKeyPrefix sets a prefix for a reward key
	RewardKeyPrefix = "rw"

	// StrategyKeyPrefix sets a prefix for a strategy key
	StrategyKeyPrefix = "st"

	// SeedKeyPrefix sets a prefix for the rng seed
	SeedKeyPrefix = "sd"
)

// RequestStoreKey returns the key to retrieve a specfic request from the store.
func RequestStoreKey(id string) []byte {
	return []byte(RequestKeyPrefix + id)
}

// RequestStoreKeyString returns the key to retrieve a specfic request from the store.
func RequestStoreKeyString(id string) string {
	return RequestKeyPrefix + id
}

// ResultStoreKey returns the key to retrieve a specfic result from the store.
func ResultStoreKey(id string) []byte {
	return []byte(ResultKeyPrefix + id)
}

// ResultStoreKeyString returns the key to retrieve a specfic result from the store.
func ResultStoreKeyString(id string) string {
	return ResultKeyPrefix + id
}

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

// RewardStoreKey returns the key to retrieve a specfic reward from the store.
func RewardStoreKey(blockHeight int64) []byte {
	//return []byte(RewardKeyPrefix + string(blockHeight))
	return []byte(RewardKeyPrefix + strconv.FormatInt(blockHeight, 10))
}

// StrategyStoreKey returns the key to retrieve a specfic strategy from the store.
func StrategyStoreKey(stratID uint64, stratName string) []byte {
	return []byte(StrategyKeyPrefix + strconv.FormatUint(stratID, 10) + stratName)
}

// SeedStoreKey returns the key to retrieve a specfic seed from the store.
func SeedStoreKey() []byte {
	return []byte(SeedKeyPrefix + "seed")
}

// RequeststoreKeyPrefixAll returns the prefix of request key (used to iterate all the requests of all requests
func RequeststoreKeyPrefixAll() []byte {
	return []byte(RequestKeyPrefix)
}
