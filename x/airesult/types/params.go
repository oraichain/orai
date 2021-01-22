package types

import (
	"fmt"

	params "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
	DefaultExpirationCount = uint64(10)
	DefaultTotalReports    = uint64(70)
)

// Parameter store keys
var (
	// TODO: Define your keys for the parameter store
	// KeyParamName          = []byte("ParamName")
	KeyExpirationCount = []byte("ExpirationCount")
	KeyTotalReports    = []byte("TotalReports")
)

// ParamKeyTable for provider module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(expirationPercentage, totalReports uint64) Params {
	return Params{
		ExpirationCount: expirationPercentage,
		TotalReports:    totalReports,
		// TODO: Create your Params Type
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
		params.NewParamSetPair(KeyExpirationCount, &p.ExpirationCount, validateExpirationCount),
		params.NewParamSetPair(KeyTotalReports, &p.TotalReports, validateTotalReports),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultExpirationCount, DefaultTotalReports)
}

func validateExpirationCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("invalid expiration count: %d", v)
	}

	return nil
}

func validateTotalReports(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("invalid total reports: %d", v)
	}

	return nil
}
