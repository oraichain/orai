package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/exported"
)

// type check for the implementation of the interface ReporterI
var _ exported.ReporterI = (*Reporter)(nil)

// Reporter is the one who send reports from the validator back to Oraichain created by a validator
type Reporter struct {
	Address   sdk.AccAddress `json:"reporter_address"`
	Name      string         `json:"reporter_name"`
	Validator sdk.ValAddress `json:"reporter_validator"`
}

// NewReporter is the constructor of the Reporter struct
func NewReporter(addr sdk.AccAddress, name string, valAddr sdk.ValAddress) Reporter {
	return Reporter{
		Address:   addr,
		Name:      name,
		Validator: valAddr,
	}
}

func (r Reporter) isEmpty() bool {
	return false
}

func (r Reporter) String() string {
	return ""
}

// GetAddress is the getter method for getting Address of a reporter
func (r Reporter) GetAddress() sdk.AccAddress {
	return r.Address
}

// GetName is the getter method for getting Name of a reporter
func (r Reporter) GetName() string {
	return r.Name
}

// GetValidator is the getter method for getting Validator of a reporter
func (r Reporter) GetValidator() sdk.ValAddress {
	return r.Validator
}
