package types

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
	DefaultOracleScriptRewardPercentage = uint64(70)
)

// // Parameter store keys
// var (
// 	// TODO: Define your keys for the parameter store
// 	// KeyParamName          = []byte("ParamName")
// 	KeyOracleScriptRewardPercentage = []byte("OracleScriptRewardPercentage")
// 	KeyMaximumCodeBytes             = []byte("MaximumCodeBytes")
// )

// // NewParams creates a new Params object
// func NewParams(rewardPercentage uint64, maxBytes uint64) Params {
// 	return Params{
// 		OracleScriptRewardPercentage: rewardPercentage,
// 		// TODO: Create your Params Type
// 		MaximumCodeBytes: maxBytes,
// 	}
// }

// // String implements the stringer interface for Params.
// func (p Params) String() string {
// 	return fmt.Sprintf(`params:
// 	OracleRewardPercentage:  %d MaximumCodeBytes: %d
// `,
// 		p.OracleScriptRewardPercentage, p.MaximumCodeBytes,
// 	)
// }

// // DefaultParams defines the parameters for this module
// func DefaultParams() Params {
// 	return NewParams(DefaultOracleScriptRewardPercentage, DefaultMaximumCodeBytes)
// }

// func validateOracleScriptRewardPercentage(i interface{}) error {
// 	v, ok := i.(uint64)
// 	if !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}

// 	if v == 0 {
// 		return fmt.Errorf("invalid oScript reward percentage: %d", v)
// 	}

// 	return nil
// }

// func validateMaximumCodeBytes(i interface{}) error {
// 	v, ok := i.(uint64)
// 	if !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}

// 	if v < 0 || v > MaximumCodeBytesThreshold {
// 		return fmt.Errorf("invalid maximum code bytes: %d", v)
// 	}

// 	return nil
// }
