package interchaintest

import (
	"context"
	"fmt"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	helpers "github.com/oraichain/orai/tests/interchaintest/helpers"
	interchaintest "github.com/strangelove-ventures/interchaintest/v4"
	"github.com/strangelove-ventures/interchaintest/v4/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v4/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v4/relayer"
	"github.com/strangelove-ventures/interchaintest/v4/testreporter"
	"github.com/strangelove-ventures/interchaintest/v4/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestJunoIBCHooks ensures the ibc-hooks middleware from osmosis works.
func TestOraiIBCHooks(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Juno and juno2
	numVals := 1
	numFullNodes := 0

	cfg2 := oraiConfig.Clone()
	cfg2.Name = "orai-counterparty"
	cfg2.ChainID = "counterparty-2"

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "orai",
			ChainConfig:   oraiConfig,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "orai-2",
			ChainConfig:   cfg2,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	const (
		path = "ibc-path"
	)

	// Get chains from the chain factory

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	client, network := interchaintest.DockerSetup(t)

	orai, orai2 := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	relayerType, relayerName := ibc.CosmosRly, "relay"

	// Get a relayer instance
	rf := interchaintest.NewBuiltinRelayerFactory(
		relayerType,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	)

	r := rf.Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(orai).
		AddChain(orai2).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  orai,
			Chain2:  orai2,
			Relayer: r,
			Path:    path,
		})

	ctx := context.Background()

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create some user accounts on both chains

	users := make([]ibc.Wallet, 0, 10)

	users = append(users,
		helpers.GetAndFundTestUserWithMnemonic(t, ctx, t.Name(), "", genesisWalletAmount, orai),
		helpers.GetAndFundTestUserWithMnemonic(t, ctx, t.Name(), "", genesisWalletAmount, orai2),
	)

	// users.app
	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, orai, orai2)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	oraiUser, orai2User := users[0], users[1]

	oraiUserAddr := oraiUser.FormattedAddress()
	orai2UserAddr := orai2User.FormattedAddress()

	// Get original account balances
	oraiOrigBal, err := orai.GetBalance(ctx, oraiUserAddr, orai.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, oraiOrigBal)

	orai2OrigBal, err := orai2.GetBalance(ctx, orai2UserAddr, orai2.Config().Denom)
	require.NoError(t, err)

	require.Equal(t, genesisWalletAmount, orai2OrigBal)
	channel, err := ibc.GetTransferChannel(ctx, r, eRep, orai.Config().ChainID, orai2.Config().ChainID)
	require.NoError(t, err)

	err = r.StartRelayer(ctx, eRep, path)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	_, contractAddr := helpers.SetupContract(t, ctx, orai2, orai2User.KeyName(), "contracts/counter.wasm", `{"count":0}`)

	// do an ibc transfer through the memo to the other chain.
	transfer := ibc.WalletAmount{
		Address: contractAddr,
		Denom:   orai.Config().Denom,
		Amount:  int64(1),
	}

	memo := ibc.TransferOptions{
		Memo: fmt.Sprintf(`{"wasm":{"contract":"%s","msg":%s}}`, contractAddr, `{"increment":{}}`),
	}

	// Send transfer for testTing
	transferTx, err := orai.SendIBCTransfer(ctx, channel.ChannelID, oraiUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)
	oraiHeight, err := orai.Height(ctx)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, orai, oraiHeight, oraiHeight+50, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 10, orai)
	require.NoError(t, err)
	// Get the IBC denom for uorai on Gaia
	oraiTokenDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, orai.Config().Denom)
	oraiIBCDenom := transfertypes.ParseDenomTrace(oraiTokenDenom).IBCDenom()

	fmt.Println("oraiIBCDenom:", oraiIBCDenom)
	orai2UserBal, err := orai2.GetBalance(ctx, contractAddr, oraiIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(1), orai2UserBal)

	// Initial transfer. Account is created by the wasm execute is not so we must do this twice to properly set up
	transferTx, err = orai.SendIBCTransfer(ctx, channel.ChannelID, oraiUserAddr, transfer, memo)
	require.NoError(t, err)
	junoHeight, err := orai.Height(ctx)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, orai, junoHeight, junoHeight+50, transferTx.Packet)
	require.NoError(t, err)
	fmt.Println("oraiIBCDenom:", oraiIBCDenom)
	orai2UserBal, err = orai2.GetBalance(ctx, contractAddr, oraiIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(2), orai2UserBal)
	//
	// // Second time, this will make the counter == 1 since the account is now created.
	// transferTx, err = orai.SendIBCTransfer(ctx, channel.ChannelID, oraiUser.KeyName(), transfer, memo)
	// require.NoError(t, err)
	// junoHeight, err = orai.Height(ctx)
	// require.NoError(t, err)
	//
	// _, err = testutil.PollForAck(ctx, orai, junoHeight-5, junoHeight+25, transferTx.Packet)
	// require.NoError(t, err)
	//
	// // Get the address on the other chain's side
	// addr := helpers.GetIBCHooksUserAddress(t, ctx, orai, channel.ChannelID, oraiUserAddr)
	// require.NotEmpty(t, addr)
	//
	// // // Get funds on the receiving chain
	// // funds := helpers.GetIBCHookTotalFunds(t, ctx, orai2, contractAddr, addr)
	// // require.Equal(t, int(1), len(funds.Data.TotalFunds))
	//
	// // var ibcDenom string
	// // for _, coin := range funds.Data.TotalFunds {
	// // 	if strings.HasPrefix(coin.Denom, "ibc/") {
	// // 		ibcDenom = coin.Denom
	// // 		break
	// // 	}
	// // }
	// // require.NotEmpty(t, ibcDenom)
	//
	// ensure the count also increased to 1 as expected.
	count := helpers.GetIBCHookCount(t, ctx, orai2, contractAddr)
	require.Equal(t, int64(1), count.Data.Count)
}
