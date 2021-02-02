package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/oraichain/orai/x/websocket/subscribe"
	"github.com/oraichain/orai/x/websocket/types"
)

const (
	flagBroadcastTimeout = "broadcast-timeout"
	flagRPCPollInterval  = "rpc-poll-interval"
	flagMaxTry           = "max-try"
	flagErrExit          = "errexit"
)

// ParseWebSocketConfig returns the default settings for WasmConfig
func ParseWebSocketConfig(flagSet *pflag.FlagSet) *subscribe.WebSocketConfig {

	cfg := subscribe.DefaultWebSocketConfig()

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

	if v, err := flagSet.GetString(flags.FlagFrom); err == nil {
		cfg.FromValidator = v
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

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdAddReporters(),
		GetCmdRemoveReporter(),
		GetCmdSubscribe(),
	)
	return txCmd
}

// GetCmdAddReporters implements the add reporters command handler.
func GetCmdAddReporters() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-reporters [reporter1] [reporter2] ...",
		Short: "Add agents authorized to submit report transactions.",
		Args:  cobra.MinimumNArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add agents authorized to submit report transactions.
Example:
$ %s tx provider add-reporters orai1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun orai1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs --from mykey
`,
				version.Name,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			validator := sdk.ValAddress(clientCtx.GetFromAddress())

			msgs := make([]sdk.Msg, len(args))
			for i, raw := range args {
				reporter, err := sdk.AccAddressFromBech32(raw)
				if err != nil {
					return err
				}
				msgs[i] = types.NewMsgAddReporter(validator, reporter, clientCtx.GetFromAddress())
				err = msgs[i].ValidateBasic()
				if err != nil {
					return err
				}
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)

		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdRemoveReporter implements the remove reporter command handler.
func GetCmdRemoveReporter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-reporter [reporter]",
		Short: "Remove an agent from the list of authorized reporters.",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Remove an agent from the list of authorized reporters.
Example:
$ %s tx provider remove-reporter band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun --from mykey
`,
				version.Name,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			validator := sdk.ValAddress(clientCtx.GetFromAddress())

			reporter, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgRemoveReporter(validator, reporter, clientCtx.GetFromAddress())
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

			cfg := ParseWebSocketConfig(cmd.Flags())
			cfg.Txf = tx.NewFactoryCLI(clientCtx, cmd.Flags())

			subscriber := subscribe.NewSubscriber(&clientCtx, cfg)
			return subscriber.Subscribe()
		},
	}

	cmd.Flags().String(flagBroadcastTimeout, "5m", "The time that the websocket will wait for tx commit")
	cmd.Flags().String(flagRPCPollInterval, "1s", "The duration of rpc poll interval")
	cmd.Flags().Uint64(flagMaxTry, 5, "The maximum number of tries to submit a report transaction")
	cmd.Flags().Bool(flagErrExit, false, "Exit on error")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
