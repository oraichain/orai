package main

import (
	"io"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/oraichain/orai/app"
	"github.com/oraichain/orai/x/websocket/websocket"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagInvCheckPeriod = "inv-check-period"
)

var invCheckPeriod uint

type CommandFuncE func(cmd *cobra.Command, args []string) error

func worker(wg *sync.WaitGroup, funcE CommandFuncE, cmd *cobra.Command, args []string, milliseconds time.Duration) {
	defer wg.Done()
	if milliseconds > 0 {
		time.Sleep(milliseconds * time.Millisecond)
	}
	funcE(cmd, args)
}

func main() {

	cdc := app.MakeCodec()
	config := sdk.GetConfig()

	// Set prefixes for the addresses of the network
	app.SetBech32AddressPrefixes(config)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "orai",
		Short:             "oracle AI Daemon (server), and custom api, websocket",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	startCmd := server.StartCmd(ctx, newApp)
	startCmd.Flags().String(flags.FlagChainID, "", "network chain ID")
	startCmd.MarkFlagRequired(flags.FlagChainID)

	// params for rest server
	flags.RegisterRestServerFlags(startCmd)

	// params for websocket
	websocket.RegisterWebSocketFlags(startCmd)

	// override RunE for start command
	runE := startCmd.RunE
	var wg sync.WaitGroup
	startCmd.RunE = func(cmd *cobra.Command, args []string) error {

		// get persistent from command param
		chainID := viper.GetString(flags.FlagChainID)

		// run start command worker
		wg.Add(1)
		go worker(&wg, runE, cmd, args, 0)

		// run rest server command worker
		// oraicli rest-server --chain-id $CHAIN_ID --laddr tcp://0.0.0.0:1317  --trust-node
		laddr := viper.GetString(flags.FlagListenAddr)
		if laddr != "" {
			restServerCmd := lcd.ServeCommand(cdc, registerRoutes)
			viper.Set(flags.FlagChainID, chainID)
			wg.Add(1)
			go worker(&wg, restServerCmd.RunE, cmd, args, 2000)
		}

		// run websocket command worker
		websocketCmd, err := websocket.ServeCommand(viper.GetString(flags.FlagHome))
		if err == nil {
			wg.Add(1)
			go worker(&wg, websocketCmd.RunE, cmd, args, 2000)
		}

		// Main: Waiting for workers to finish
		wg.Wait()

		return nil
	}

	rootCmd.AddCommand(
		startCmd,
		flags.LineBreak,
		server.VersionCmd(ctx),
	)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "ORAI", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	return app.NewOraichainApp(
		logger, db, traceStore, true, invCheckPeriod,
		viper.GetString(flags.FlagHome),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}
