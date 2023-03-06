package e2e

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testBankTokenTransfer() {
	s.Run("send_photon_between_accounts", func() {
		var err error
		senderAddress := s.chainA.validators[0].keyInfo.GetAddress()
		sender := senderAddress.String()

		recipientAddress := s.chainA.validators[1].keyInfo.GetAddress()
		recipient := recipientAddress.String()

		chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

		var (
			beforeSenderuoraiBalance    sdk.Coin
			beforeRecipientuoraiBalance sdk.Coin
		)

		s.Require().Eventually(
			func() bool {
				beforeSenderuoraiBalance, err = getSpecificBalance(chainAAPIEndpoint, sender, uoraiDenom)
				s.Require().NoError(err)

				beforeRecipientuoraiBalance, err = getSpecificBalance(chainAAPIEndpoint, recipient, uoraiDenom)
				s.Require().NoError(err)

				return beforeSenderuoraiBalance.IsValid() && beforeRecipientuoraiBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		s.execBankSend(s.chainA, 0, sender, recipient, tokenAmount.String(), standardFees.String(), false)

		s.Require().Eventually(
			func() bool {
				afterSenderuoraiBalance, err := getSpecificBalance(chainAAPIEndpoint, sender, uoraiDenom)
				s.Require().NoError(err)

				afterRecipientuoraiBalance, err := getSpecificBalance(chainAAPIEndpoint, recipient, uoraiDenom)
				s.Require().NoError(err)

				decremented := beforeSenderuoraiBalance.Sub(tokenAmount).Sub(standardFees).IsEqual(afterSenderuoraiBalance)
				incremented := beforeRecipientuoraiBalance.Add(tokenAmount).IsEqual(afterRecipientuoraiBalance)

				return decremented && incremented
			},
			time.Minute,
			5*time.Second,
		)
	})
}
