package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/exported"
)

// type check for the implementation of the interface ValidatorI
var _ exported.ValidatorI = (*Validator)(nil)

// Validator mimics the original validator to store information of those that execute the oScript
type Validator struct {
	Address     sdk.ValAddress `json:"address"`
	VotingPower int64          `json:"voting_power"`
	Status      string         `json:"status"`
}

// NewValidator is the constructor of the validator struct
func NewValidator(
	address sdk.ValAddress,
	votingPower int64,
	status string,
) Validator {
	return Validator{
		Address:     address,
		VotingPower: votingPower,
		Status:      status,
	}
}

func (v Validator) isEmpty() bool {
	return false
}

// GetAddress is getter method for Validator struct
func (v Validator) GetAddress() sdk.ValAddress {
	return v.Address
}

// GetVotingPower is getter method for Validator struct
func (v Validator) GetVotingPower() int64 {
	return v.VotingPower
}

// GetStatus is getter method for Validator struct
func (v Validator) GetStatus() string {
	return v.Status
}

// type check for the implementation of the interface ValidatorI
var _ exported.ValResultI = (*ValResult)(nil)

// ValResult stores the result information from a validator that has executed the oracle script
type ValResult struct {
	Validator    exported.ValidatorI `json:"validator"`
	Result       []byte              `json:"result"`
	ResultStatus string              `json:"result_status"` // fail or success
}

// NewValResult is a constructor for the validator result
func NewValResult(
	val exported.ValidatorI,
	result []byte,
	status string,
) ValResult {
	return ValResult{
		Validator:    val,
		Result:       result,
		ResultStatus: status,
	}
}

// GetValidator is getter method for ValResult struct
func (valRes ValResult) GetValidator() exported.ValidatorI {
	return valRes.Validator
}

// GetResult is getter method for ValResult struct
func (valRes ValResult) GetResult() []byte {
	return valRes.Result
}

// GetResult is getter method for ValResult struct
func (valRes ValResult) GetResultStatus() string {
	return valRes.ResultStatus
}
