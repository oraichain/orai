package cli

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/subscribe"
	"github.com/oraichain/orai/x/aioracle/types"
	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tendermint/tendermint/libs/log"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "aioracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdSetAIOracle(),
		GetCmdSubscribe(),
	)
	return txCmd
}

const (
	flagBroadcastTimeout = "broadcast-timeout"
	flagRPCPollInterval  = "rpc-poll-interval"
	flagMaxTry           = "max-try"
	flagErrExit          = "errexit"
)

// GetCmdSetAIOracle is the CLI command for sending a SetAIOracle transaction
func GetCmdSetAIOracle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-orai [oracle-script-name] [input] [transaction-fees] [validator-count]",
		Short: "Set a new ai oracle request and set result into the system",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valCount, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetAIOracleReq(ksuid.New().String(), contractAddr, clientCtx.GetFromAddress(), args[3], int64(valCount), []byte(args[1]))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ParseaioracleConfig returns the default settings for WasmConfig
func ParseaioracleConfig(flagSet *pflag.FlagSet) *subscribe.AIOracleConfig {

	cfg := subscribe.DefaultAIOracleConfig()

	if v, err := flagSet.GetString(flagBroadcastTimeout); err == nil {
		if cfg.BroadcastTimeout, err = time.ParseDuration(v); err != nil {
			return cfg
		}
	}
	if v, err := flagSet.GetString(flagRPCPollInterval); err == nil {
		if cfg.RPCPollInterval, err = time.ParseDuration(v); err != nil {
			return cfg
		}
	}

	if v, err := flagSet.GetUint64(flagMaxTry); err == nil {
		cfg.MaxTry = v
	}

	if v, err := flagSet.GetBool(flagErrExit); err == nil {
		cfg.ErrExit = v
	}

	if v, err := flagSet.GetString(flags.FlagLogLevel); err == nil {
		if cfg.AllowLogLevel, err = log.AllowLevel(v); err != nil {
			return cfg
		}
	}

	if v, err := flagSet.GetString(flags.FlagLogLevel); err == nil {
		if cfg.AllowLogLevel, err = log.AllowLevel(v); err != nil {
			return cfg
		}
	}

	if v, err := flagSet.GetString(flags.FlagFees); err == nil {
		if cfg.Fees, err = sdk.ParseCoinsNormalized(v); err != nil {
			return cfg
		}
	}

	return cfg
}

// GetCmdSubscribe implements the subscribe handler.
func GetCmdSubscribe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe to AI request log to submit report transactions.",
		RunE: func(cmd *cobra.Command, args []string) error {

			// txContext no require for block height, it just broadcast tx
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cfg := ParseaioracleConfig(cmd.Flags())
			cfg.Txf = tx.NewFactoryCLI(clientCtx, cmd.Flags())

			subscriber := subscribe.NewSubscriber(&clientCtx, cfg)
			return subscriber.Subscribe()
		},
	}

	cmd.Flags().String(flagBroadcastTimeout, "5m", "The time that the aioracle will wait for tx commit")
	cmd.Flags().String(flagRPCPollInterval, "1s", "The duration of rpc poll interval")
	cmd.Flags().Uint64(flagMaxTry, 5, "The maximum number of tries to submit a report transaction")
	cmd.Flags().Bool(flagErrExit, false, "Exit on error")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
