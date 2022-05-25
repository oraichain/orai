package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
	MaximumRequestBytesThreshold = 1 * 1024 * 1024 // 1MB
)

// Parameter store keys
var (
	// TODO: Define your keys for the parameter store
	// KeyParamName          = []byte("ParamName")
	KeyMaximumRequestBytes = []byte("MaximumRequestBytes")
)

// ParamKeyTable for provider module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(maximumReqBytes uint64) Params {
	// if the value exceeds the threshold => default is the maximum value
	if maximumReqBytes > MaximumRequestBytesThreshold {
		return Params{
			MaximumRequestBytes: MaximumRequestBytesThreshold,
		}
	}
	return Params{
		// TODO: Create your Params Type
		MaximumRequestBytes: maximumReqBytes,
	}
}

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`params:  MaximumRequestBytes: %d`, p.MaximumRequestBytes)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
		paramtypes.NewParamSetPair(KeyMaximumRequestBytes, &p.MaximumRequestBytes, validateMaximumRequestBytes),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(MaximumRequestBytesThreshold)
}

func validateMaximumRequestBytes(i interface{}) error {
	v, ok := i.(int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("invalid maximum bytes: %d", v)
	}

	return nil
}
