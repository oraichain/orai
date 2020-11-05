package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/exported"
)

// Implements Test Case interface
var _ exported.TestCaseI = (*TestCase)(nil)

// TestCase stores the test case of a request to an AI model
type TestCase struct {
	Owner       sdk.AccAddress `json:"owner"`
	Name        string         `json:"name"`
	Fees        sdk.Coins      `json:"transaction_fee"`
	Description string         `json:"description"`
}

// implement fmt.Stringer
func (tc TestCase) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
Owner: %s Name: %s Fees: %s Description: %s`, tc.Owner, tc.Name, tc.Fees.String(), tc.Description))
}

// NewTestCase is the constructor of the test case struct
func NewTestCase(
	name string,
	owner sdk.AccAddress,
	fees sdk.Coins,
	des string,
) TestCase {
	return TestCase{
		Name:        name,
		Owner:       owner,
		Fees:        fees,
		Description: des,
	}
}

// SetName is the setter function for updating the TestCase's name
func (tc TestCase) SetName(newName string) error {
	if len(newName) == 0 {
		return errors.New("Cannot set a new name as an empty string")
	}
	tc.Name = newName
	return nil
}

// GetName is the getter function for getting the TestCase's name
func (tc TestCase) GetName() string {
	return tc.Name
}

// SetDescription is the setter function for updating the TestCase's description
func (tc TestCase) SetDescription(newDescription string) error {
	if len(newDescription) == 0 {
		return errors.New("Cannot set a new description as an empty string")
	}
	tc.Description = newDescription
	return nil
}

// GetDescription is the getter function for getting the TestCase's description
func (tc TestCase) GetDescription() string {
	return tc.Description
}

// SetOwner is the setter function for updating the TestCase's owner
func (tc TestCase) SetOwner(newOwner sdk.AccAddress) error {
	if newOwner.Empty() {
		return errors.New("Cannot set a new owner as empty")
	}
	tc.Owner = newOwner
	return nil
}

// GetOwner is the getter function for getting the TestCase's owner
func (tc TestCase) GetOwner() sdk.AccAddress {
	return tc.Owner
}

// SetFees is the setter function for updating the TestCase's fees
func (tc TestCase) SetFees(fees sdk.Coins) error {
	if fees.Empty() {
		return errors.New("Cannot set a new fees as empty")
	}
	tc.Fees = fees
	return nil
}

// GetFees is the getter function for getting the TestCase's fees
func (tc TestCase) GetFees() sdk.Coins {
	return tc.Fees
}
