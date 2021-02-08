package types

import "strconv"

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
	DefaultQueryPage = 1

	// DefaultQueryLimit sets the default query limit value
	DefaultQueryLimit = 5

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"

	// RequestStatusPending is the status pending of the request
	RequestStatusPending = "pending"

	// RequestStatusFinished is the status finished of the request
	RequestStatusFinished = "finished"

	// RequestStatusExpired is the status expired of the request
	RequestStatusExpired = "expired"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (

	// ResultKeyPrefix sets a prefix for a result key
	ResultKeyPrefix = "res"

	// RewardKeyPrefix sets a prefix for a reward key
	RewardKeyPrefix = "rw"
)

// ResultStoreKey returns the key to retrieve a specfic result from the store.
func ResultStoreKey(id string) []byte {
	return []byte(ResultKeyPrefix + id)
}

// RewardStoreKey returns the key to retrieve a specfic reward from the store.
func RewardStoreKey(blockHeight int64) []byte {
	//return []byte(RewardKeyPrefix + string(blockHeight))
	return []byte(RewardKeyPrefix + strconv.FormatInt(blockHeight, 10))
}
