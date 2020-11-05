package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oraichain/orai/x/provider/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Implements Data Source interface
var _ exported.AIDataSourceI = (*AIDataSource)(nil)

// AIDataSource is a struct for storing AIDataSource information of a provider
type AIDataSource struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        sdk.Coins      `json:"transaction_fee"`
	Description string         `json:"description"`
}

// implement fmt.Stringer
func (ds AIDataSource) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s
Owner: %s Fees: %s Description: %s`, ds.Name, ds.Owner, ds.Fees.String(), ds.Description))
}

// NewAIDataSource is the constructor of the data source struct
func NewAIDataSource(
	name string,
	owner sdk.AccAddress,
	fees sdk.Coins,
	des string,
) AIDataSource {
	return AIDataSource{
		Name:        name,
		Owner:       owner,
		Fees:        fees,
		Description: des,
	}
}

// SetName is the setter function for updating the AIDataSource's name
func (ds AIDataSource) SetName(newName string) error {
	if len(newName) == 0 {
		return errors.New("Cannot set a new name as an empty string")
	}
	ds.Name = newName
	return nil
}

// GetName is the getter function for getting the AIDataSource's name
func (ds AIDataSource) GetName() string {
	return ds.Name
}

// SetDescription is the setter function for updating the AIDataSource's description
func (ds AIDataSource) SetDescription(newDescription string) error {
	if len(newDescription) == 0 {
		return errors.New("Cannot set a new description as an empty string")
	}
	ds.Description = newDescription
	return nil
}

// GetDescription is the getter function for getting the AIDataSource's description
func (ds AIDataSource) GetDescription() string {
	return ds.Description
}

// SetOwner is the setter function for updating the AIDataSource's owner
func (ds AIDataSource) SetOwner(newOwner sdk.AccAddress) error {
	if newOwner.Empty() {
		return errors.New("Cannot set a new owner as empty")
	}
	ds.Owner = newOwner
	return nil
}

// GetOwner is the getter function for getting the AIDataSource's owner
func (ds AIDataSource) GetOwner() sdk.AccAddress {
	return ds.Owner
}

// SetFees is the setter function for updating the AIDataSource's fees
func (ds AIDataSource) SetFees(fees sdk.Coins) error {
	if fees.Empty() {
		return errors.New("Cannot set a new fees as empty")
	}
	ds.Fees = fees
	return nil
}

// GetFees is the getter function for getting the AIDataSource's fees
func (ds AIDataSource) GetFees() sdk.Coins {
	return ds.Fees
}
