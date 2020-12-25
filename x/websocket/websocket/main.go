package websocket

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagValidator        = "validator"
	flagLogLevel         = "log-level"
	flagExecutor         = "executor"
	flagBroadcastTimeout = "broadcast-timeout"
	flagRPCPollInterval  = "rpc-poll-interval"
	flagMaxTry           = "max-try"
	// Bech32MainPrefix is the prefix of different addresses
	Bech32MainPrefix = "orai"

	defaultFees = "5000orai"
)

// Config data structure for the websocket daemon.
type Config struct {
	ChainID          string `mapstructure:"chain-id"`          // ChainID of the target chain
	NodeURI          string `mapstructure:"node"`              // Remote RPC URI of Oraichain node to connect to
	Validator        string `mapstructure:"validator"`         // The validator address that I'm responsible for
	GasPrices        string `mapstructure:"gas-prices"`        // Gas prices of the transaction
	GasAdjustment    string `mapstructure:"gas-adjustment"`    // Gas adjustment from the gas
	Gas              string `mapstructure:"gas"`               // gas used for the transaction
	LogLevel         string `mapstructure:"log-level"`         // Log level of the logger
	Executor         string `mapstructure:"executor"`          // Executor name and URL (example: "Executor name:URL")
	BroadcastTimeout string `mapstructure:"broadcast-timeout"` // The time that the websocket will wait for tx commit
	RPCPollInterval  string `mapstructure:"rpc-poll-interval"` // The duration of rpc poll interval
	MaxTry           uint64 `mapstructure:"max-try"`           // The maximum number of tries to submit a report transaction
}

// Global instances.
var (
	cfg     Config
	keybase keys.Keybase
)

func restoreConfig(home string) error {
	viper.SetConfigFile(path.Join(home, "websocket.yaml"))
	_ = viper.ReadInConfig() // If we fail to read config file, we'll just rely on cmd flags.
	return viper.Unmarshal(&cfg)
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}

	err = restoreConfig(home)

	if err != nil {
		return err
	}
	return nil
}

// Main is the entry point of the validator websocket
func Main() {

	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	accountPrefix := Bech32MainPrefix
	validatorPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	consensusPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	config.SetBech32PrefixForAccount(accountPrefix, accountPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(validatorPrefix, validatorPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(consensusPrefix, consensusPrefix+sdk.PrefixPublic)
	config.Seal()

	ctx := &Context{}
	rootCmd := &cobra.Command{
		Use:   "websocket",
		Short: "Orai websocket to subscribe and response to AI requests",
	}

	rootCmd.AddCommand(configCmd(), keysCmd(ctx), runCmd(ctx), version.Cmd)
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		home, err := rootCmd.PersistentFlags().GetString(flags.FlagHome)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(home, os.ModePerm); err != nil {
			return err
		}
		keybase, err = keys.NewKeyring("orai", "test", home, nil)
		if err != nil {
			return err
		}
		return initConfig(rootCmd)
	}
	rootCmd.PersistentFlags().String(flags.FlagHome, os.ExpandEnv("$PWD/.websocket"), "home directory")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ServeCommand supports orai command to have a quick start for a node
func ServeCommand(home string) (cmd *cobra.Command, err error) {
	// Configure cobra to sort commands
	ctx := &Context{}
	keybase, err = keys.NewKeyring("orai", "test", home, nil)
	err = restoreConfig(home)
	cmd = runCmd(ctx)
	return cmd, err
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config [key] [value]",
		Aliases: []string{"c"},
		Short:   "Set websocket configuration environment",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.Set(args[0], args[1])
			return viper.WriteConfig()
		},
	}
	return cmd
}
