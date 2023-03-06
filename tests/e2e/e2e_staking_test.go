package e2e

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testStaking() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorB := s.chainA.validators[1]
	validatorAAddr := validatorA.keyInfo.GetAddress()
	validatorBAddr := validatorB.keyInfo.GetAddress()

	validatorAddressA := sdk.ValAddress(validatorAAddr).String()
	validatorAddressB := sdk.ValAddress(validatorBAddr).String()

	delegatorAddress := s.chainA.genesisAccounts[2].keyInfo.GetAddress().String()

	fees := sdk.NewCoin(uoraiDenom, sdk.NewInt(10))

	delegationAmount := sdk.NewInt(500000000)
	delegation := sdk.NewCoin(uoraiDenom, delegationAmount) // 500 atom

	// Alice delegate uorai to Validator A
	s.executeDelegate(s.chainA, 0, delegation.String(), validatorAddressA, delegatorAddress, oraiHomePath, fees.String())

	// Validate delegation successful
	s.Require().Eventually(
		func() bool {
			res, err := queryDelegation(chainEndpoint, validatorAddressA, delegatorAddress)
			amt := res.GetDelegationResponse().GetDelegation().GetShares()
			s.Require().NoError(err)

			return amt.Equal(sdk.NewDecFromInt(delegationAmount))
		},
		20*time.Second,
		5*time.Second,
	)

	// Alice re-delegate uorai from Validator A to Validator B
	s.executeRedelegate(s.chainA, 0, delegation.String(), validatorAddressA, validatorAddressB, delegatorAddress, oraiHomePath, fees.String())

	// Validate re-delegation successful
	s.Require().Eventually(
		func() bool {
			res, err := queryDelegation(chainEndpoint, validatorAddressB, delegatorAddress)
			amt := res.GetDelegationResponse().GetDelegation().GetShares()
			s.Require().NoError(err)

			return amt.Equal(sdk.NewDecFromInt(delegationAmount))
		},
		20*time.Second,
		5*time.Second,
	)
}
