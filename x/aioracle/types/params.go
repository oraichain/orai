package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
	DefaultOracleReqBytesThreshold = 1 * 1024 * 1024 // 1MB
	DefaultOracleResBytesThreshold = 1 * 1024 * 1024 // 1MB
	DefaultOracleRewardPercentages = 40
	DefaultReportPercentages       = uint64(70)
)

// Parameter store keys
var (
	// TODO: Define your keys for the parameter store
	// KeyParamName          = []byte("ParamName")
	KeyMaximumAIOracleReqBytes   = []byte("MaximumRequestBytes")
	KeyMaximumAIOracleResBytes   = []byte("RewardPercentageBytes")
	KeyAIOracleRewardPercentages = []byte("KeyAIOracleRewardPercentages")
	KeyReportPercentages         = []byte("KeyReportPercentages")
)

// ParamKeyTable for provider module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(maximumReqBytes, maximumResBytes, rewardPercentages, reportsPercentages uint64) Params {
	return Params{
		// TODO: Create your Params Type
		MaximumAiOracleRequestBytes:  maximumReqBytes,
		MaximumAiOracleResponseBytes: maximumResBytes,
		RewardAiOraclePercentages:    rewardPercentages,
		ReportsPercentages:           reportsPercentages,
	}
}

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`params:  MaximumAiOracleRequestBytes: %d MaximumAiOracleRequestBytes: %d RewardAiOraclePercentages: %d ReportsPercentages: %d`, p.MaximumAiOracleRequestBytes, p.MaximumAiOracleRequestBytes, p.RewardAiOraclePercentages, p.ReportsPercentages)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
		paramtypes.NewParamSetPair(KeyMaximumAIOracleReqBytes, &p.MaximumAiOracleRequestBytes, validateMaximumAiOracleReqBytes),
		paramtypes.NewParamSetPair(KeyMaximumAIOracleResBytes, &p.MaximumAiOracleResponseBytes, validateMaximumAiOracleResBytes),
		paramtypes.NewParamSetPair(KeyAIOracleRewardPercentages, &p.RewardAiOraclePercentages, validateAIOracleRewardPercentages),
		paramtypes.NewParamSetPair(KeyReportPercentages, &p.ReportsPercentages, validateReportPercentages),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultOracleReqBytesThreshold, DefaultOracleResBytesThreshold, DefaultOracleRewardPercentages, DefaultReportPercentages)
}

func validateMaximumAiOracleReqBytes(i interface{}) error {
	v, ok := i.(int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("invalid maximum bytes: %d", v)
	}

	return nil
}

func validateMaximumAiOracleResBytes(i interface{}) error {
	v, ok := i.(int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("invalid maximum bytes: %d", v)
	}

	return nil
}
func validateAIOracleRewardPercentages(i interface{}) error {
	v, ok := i.(int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("invalid maximum bytes: %d", v)
	}

	return nil
}
func validateReportPercentages(i interface{}) error {
	v, ok := i.(int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("invalid maximum bytes: %d", v)
	}

	return nil
}
