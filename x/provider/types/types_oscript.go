package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oraichain/orai/x/provider/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Implements Oracle Script interface
var _ exported.OracleScriptI = (*OracleScript)(nil)

// OracleScript is a struct for storing oracle script information
type OracleScript struct {
	Name        string         `json:"name"`
	Owner       sdk.AccAddress `json:"owner"`
	Description string         `json:"description"`
	MinimumFees sdk.Coins      `json:"minimum_fees"`
	DSources    []string       `json:"data_sources"`
	TCases      []string       `json:"test_cases"`
}

// implement fmt.Stringer
func (os OracleScript) String() string {
	dsListNames := strings.Join(os.DSources[:], "\n")
	tcListNames := strings.Join(os.TCases[:], "\n")
	return strings.TrimSpace(fmt.Sprintf(`Name: %s
Owner: %s Description: %s Minimum Fees: %s Data Sources: %s Test Cases: %s`, os.Name, os.Owner, os.Description, os.MinimumFees, dsListNames, tcListNames))
}

// NewOracleScript is the constructor of the oScript struct
func NewOracleScript(
	name string,
	owner sdk.AccAddress,
	des string,
	minimumFees sdk.Coins,
	dSources []string,
	tCases []string,
) OracleScript {
	return OracleScript{
		Name:        name,
		Owner:       owner,
		Description: des,
		MinimumFees: minimumFees,
		DSources:    dSources,
		TCases:      tCases,
	}
}

// SetName is the setter function for updating the OracleScript's name
func (os OracleScript) SetName(newName string) error {
	if len(newName) == 0 {
		return errors.New("Cannot set a new name as an empty string")
	}
	os.Name = newName
	return nil
}

// GetName is the getter function for getting the OracleScript's name
func (os OracleScript) GetName() string {
	return os.Name
}

// SetDescription is the setter function for updating the OracleScript's description
func (os OracleScript) SetDescription(newDescription string) error {
	if len(newDescription) == 0 {
		return errors.New("Cannot set a new description as an empty string")
	}
	os.Description = newDescription
	return nil
}

// GetDescription is the getter function for getting the OracleScript's description
func (os OracleScript) GetDescription() string {
	return os.Description
}

// SetOwner is the setter function for updating the OracleScript's owner
func (os OracleScript) SetOwner(newOwner sdk.AccAddress) error {
	if newOwner.Empty() {
		return errors.New("Cannot set a new owner as empty")
	}
	os.Owner = newOwner
	return nil
}

// GetOwner is the getter function for getting the OracleScript's owner
func (os OracleScript) GetOwner() sdk.AccAddress {
	return os.Owner
}

// SetMinimumFees is the setter function for updating the OracleScript's fees
func (os OracleScript) SetMinimumFees(minFees sdk.Coins) error {
	if minFees.Empty() {
		return errors.New("Cannot set a new fees as empty")
	}
	os.MinimumFees = minFees
	return nil
}

// GetMinimumFees is the getter function for getting the OracleScript's fees
func (os OracleScript) GetMinimumFees() sdk.Coins {
	return os.MinimumFees
}

// SetDSources is the setter function for updating the OracleScript's data sources
func (os OracleScript) SetDSources(dSources []string) error {
	if len(dSources) == 0 {
		return errors.New("The data source list is empty")
	}
	os.DSources = dSources
	return nil
}

// GetDSources is the getter function for getting the OracleScript's data sources
func (os OracleScript) GetDSources() []string {
	return os.DSources
}

// SetTCases is the setter function for updating the OracleScript's test cases
func (os OracleScript) SetTCases(tCases []string) error {
	if len(tCases) == 0 {
		return errors.New("The test case list is empty")
	}
	os.TCases = tCases
	return nil
}

// GetTCases is the getter function for getting the OracleScript's test cases
func (os OracleScript) GetTCases() []string {
	return os.TCases
}
