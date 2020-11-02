package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
	DefaultOracleScriptRewardPercentage = uint64(70)
	DefaultExpirationCount              = uint64(10)
)

// Parameter store keys
var (
	// TODO: Define your keys for the parameter store
	// KeyParamName          = []byte("ParamName")
	KeyOracleScriptRewardPercentage = []byte("OracleScriptRewardPercentage")
	KeyExpirationCount              = []byte("ExpirationCount")
)

// ParamKeyTable for provider module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for provider at genesis
type Params struct {
	// TODO: Add your Paramaters to the Paramter struct
	// KeyParamName string `json:"key_param_name"`
	OracleScriptRewardPercentage uint64 `json:"oscript_reward_percentage"`
	ExpirationCount              uint64 `json:"expiration_count"`
}

// NewParams creates a new Params object
func NewParams(rewardPercentage uint64, expirationPercentage uint64) Params {
	return Params{
		OracleScriptRewardPercentage: rewardPercentage,
		ExpirationCount:              expirationPercentage,
		// TODO: Create your Params Type
	}
}

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`params:
	OracleRewardPercentage:  %d
	ExpirationCount: %d
`,
		p.OracleScriptRewardPercentage,
		p.ExpirationCount,
	)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
		params.NewParamSetPair(KeyOracleScriptRewardPercentage, &p.OracleScriptRewardPercentage, validateOracleScriptRewardPercentage),
		params.NewParamSetPair(KeyExpirationCount, &p.ExpirationCount, validateExpirationCount),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultOracleScriptRewardPercentage, DefaultExpirationCount)
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

func validateExpirationCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid expiration count: %d", v)
	}

	return nil
}
