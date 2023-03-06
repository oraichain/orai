package e2e

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testDistribution() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorB := s.chainA.validators[1]
	validatorBAddr := validatorB.keyInfo.GetAddress()

	valOperAddressA := sdk.ValAddress(validatorBAddr).String()

	delegatorAddress := s.chainA.genesisAccounts[2].keyInfo.GetAddress().String()

	newWithdrawalAddress := s.chainA.genesisAccounts[3].keyInfo.GetAddress().String()
	fees := sdk.NewCoin(uoraiDenom, sdk.NewInt(1000))

	beforeBalance, err := getSpecificBalance(chainEndpoint, newWithdrawalAddress, uoraiDenom)
	s.Require().NoError(err)
	if beforeBalance.IsNil() {
		beforeBalance = sdk.NewCoin(uoraiDenom, sdk.NewInt(0))
	}

	s.execSetWithdrawAddress(s.chainA, 0, fees.String(), delegatorAddress, newWithdrawalAddress, oraiHomePath)

	// Verify
	s.Require().Eventually(
		func() bool {
			res, err := queryDelegatorWithdrawalAddress(chainEndpoint, delegatorAddress)
			s.Require().NoError(err)

			return res.WithdrawAddress == newWithdrawalAddress
		},
		10*time.Second,
		5*time.Second,
	)

	s.execWithdrawReward(s.chainA, 0, delegatorAddress, valOperAddressA, oraiHomePath)
	s.Require().Eventually(
		func() bool {
			afterBalance, err := getSpecificBalance(chainEndpoint, newWithdrawalAddress, uoraiDenom)
			s.Require().NoError(err)

			return afterBalance.IsGTE(beforeBalance)
		},
		10*time.Second,
		5*time.Second,
	)
}

/*
fundCommunityPool tests the funding of the community pool on behalf of the distribution module.
Test Benchmarks:
1. Validation that balance of the distribution module account before funding
2. Execution funding the community pool
3. Verification that correct funds have been deposited to distribution module account
*/
func (s *IntegrationTestSuite) fundCommunityPool() {
	chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	sender := s.chainA.validators[0].keyInfo.GetAddress()

	beforeDistuoraiBalance, _ := getSpecificBalance(chainAAPIEndpoint, distModuleAddress, tokenAmount.Denom)
	if beforeDistuoraiBalance.IsNil() {
		// Set balance to 0 if previous balance does not exist
		beforeDistuoraiBalance = sdk.NewInt64Coin(uoraiDenom, 0)
	}

	s.execDistributionFundCommunityPool(s.chainA, 0, sender.String(), tokenAmount.String(), standardFees.String())

	// there are still tokens being added to the community pool through block production rewards but they should be less than 500 tokens
	marginOfErrorForBlockReward := sdk.NewInt64Coin(uoraiDenom, 500)

	s.Require().Eventually(
		func() bool {
			afterDistPhotonBalance, err := getSpecificBalance(chainAAPIEndpoint, distModuleAddress, tokenAmount.Denom)
			s.Require().NoErrorf(err, "Error getting balance: %s", afterDistPhotonBalance)

			return afterDistPhotonBalance.Sub(beforeDistuoraiBalance.Add(tokenAmount.Add(standardFees))).IsLT(marginOfErrorForBlockReward)
		},
		15*time.Second,
		5*time.Second,
	)
}
