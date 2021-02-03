package types

const (
	// ModuleName is the name of the module
	ModuleName = "airequest"

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
