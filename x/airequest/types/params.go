package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
	DefaultMaximumRequestBytes   = 1024
	MaximumRequestBytesThreshold = 1 * 1024 * 1024 // 1MB
)

// Parameter store keys
var (
	// TODO: Define your keys for the parameter store
	// KeyParamName          = []byte("ParamName")
	KeyMaximumRequestBytes = []byte("MaximumRequestBytes")
)

// ParamKeyTable for provider module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for provider at genesis
type Params struct {
	// TODO: Add your Paramaters to the Paramter struct
	// KeyParamName string `json:"key_param_name"`
	MaximumRequestBytes uint64 `json:"maximum_request_bytes"`
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
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
		params.NewParamSetPair(KeyMaximumRequestBytes, &p.MaximumRequestBytes, validateMaximumRequestBytes),
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

	if v < 0 || v > MaximumRequestBytesThreshold {
		return fmt.Errorf("invalid maximum bytes: %d", v)
	}

	return nil
}
