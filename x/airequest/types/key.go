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

	// IPFSUrl is he default URL of ipfs gateway to store user images
	IPFSUrl = "http://164.90.180.95:5001/api/v0"

	//IPFSAdd is the path for adding a file onto IPFS
	IPFSAdd = "/add"

	//IPFSCat is the path for retrieving a file from IPFS into the system
	IPFSCat = "/cat"

	// DefaultQueryPage sets the default page query value
	DefaultQueryPage = "1"

	// DefaultQueryLimit sets the default query limit value
	DefaultQueryLimit = "5"

	// RngSeedSize is the size of the random seed for validator sampling.
	RngSeedSize = 64

	// NumSeedRemoval is the number of bytes an old seed is removed to generate a new one
	NumSeedRemoval = 1
	// Denom is the denominator of the currency
	Denom = "orai"
)

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
