package exported

// // These interfaces will be used by other modules when they need to interact with the attributes of the struct

// // ProviderI is the generic Provider interface that applies for all types of providers
// type ProviderI interface {
// 	GetName() string
// 	GetDescription() string
// 	GetOwner() sdk.AccAddress
// }

// // OracleScriptI expected oracle script functions
// type OracleScriptI interface {
// 	ProviderI
// 	GetMinimumFees() sdk.Coins
// 	GetDSources() []string
// 	GetTCases() []string
// }

// // AIDataSourceI expected data source functions
// type AIDataSourceI interface {
// 	ProviderI
// 	GetFees() sdk.Coins
// }

// // TestCaseI expected test case functions
// type TestCaseI interface {
// 	ProviderI
// 	GetFees() sdk.Coins
// }
