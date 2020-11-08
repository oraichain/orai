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
