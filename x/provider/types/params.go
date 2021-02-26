package types

import (
	"fmt"

	params "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
	DefaultOracleScriptRewardPercentage = uint64(60)
	DefaultMaximumCodeBytes             = uint64(1024 * 1024) // 1MB
	MaximumCodeBytesThreshold           = 4 * 1024 * 1024     // 4MB
)

// Parameter store keys
var (
	// TODO: Define your keys for the parameter store
	// KeyParamName          = []byte("ParamName")
	KeyOracleScriptRewardPercentage = []byte("OracleScriptRewardPercentage")
	KeyMaximumCodeBytes             = []byte("MaximumCodeBytes")
)

// ParamKeyTable for provider module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(rewardPercentage uint64, maxBytes uint64) Params {
	return Params{
		OracleScriptRewardPercentage: rewardPercentage,
		// TODO: Create your Params Type
		MaximumCodeBytes: maxBytes,
	}
}

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`params:
	OracleRewardPercentage:  %d MaximumCodeBytes: %d
`,
		p.OracleScriptRewardPercentage, p.MaximumCodeBytes,
	)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
		params.NewParamSetPair(KeyOracleScriptRewardPercentage, &p.OracleScriptRewardPercentage, validateOracleScriptRewardPercentage),
		params.NewParamSetPair(KeyMaximumCodeBytes, &p.MaximumCodeBytes, validateMaximumCodeBytes),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultOracleScriptRewardPercentage, DefaultMaximumCodeBytes)
}

func validateOracleScriptRewardPercentage(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid oScript reward percentage: %d", v)
	}

	return nil
}

func validateMaximumCodeBytes(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 || v > MaximumCodeBytesThreshold {
		return fmt.Errorf("invalid maximum code bytes: %d", v)
	}

	return nil
}
