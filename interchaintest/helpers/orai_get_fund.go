package helpers

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v4/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v4/ibc"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

var chars = []byte("abcdefghijklmnopqrstuvwxyz")

const (
	FaucetAccountKeyName = "faucet"
)

func RandLowerCaseLetterString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func GetAndFundTestUserWithMnemonic(
	t *testing.T,
	ctx context.Context,
	keyNamePrefix, mnemonic string,
	amount int64,
	chain ibc.Chain,
) ibc.Wallet {
	var eg errgroup.Group
	chainCfg := chain.Config()
	keyName := fmt.Sprintf("%s-%s-%s", keyNamePrefix, chainCfg.ChainID, RandLowerCaseLetterString(3))
	user, err := chain.BuildWallet(ctx, keyName, mnemonic)
	require.NoError(t, err)
	cosmosChain := chain.(*cosmos.CosmosChain)
	var cn *cosmos.ChainNode

	if len(cosmosChain.FullNodes) > 0 {
		cn = cosmosChain.FullNodes[0]
	} else {
		cn = cosmosChain.Validators[0]
	}

	eg.Go(func() error {
		_, err = cn.ExecTx(ctx,
			FaucetAccountKeyName,
			"send", FaucetAccountKeyName, user.FormattedAddress(),
			fmt.Sprintf("%d%s", amount, chainCfg.Denom),
		)
		if err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, eg.Wait())
	return user
}
