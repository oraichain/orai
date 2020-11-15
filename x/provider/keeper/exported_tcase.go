package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/exported"
)

// This file implements all the exported functions of the Test Case struct that are exported by the Keeper

// GetTestCaseI returns the test case object given the name of the test case
func (k Keeper) GetTestCaseI(ctx sdk.Context, name string) (exported.TestCaseI, error) {
	testCase, err := k.GetTestCase(ctx, name)
	if err != nil {
		return nil, err
	}
	return testCase, nil
}

// DefaultTestCaseI returns the default test case object
func (k Keeper) DefaultTestCaseI() exported.TestCaseI {
	return k.DefaultTestCase()
}
