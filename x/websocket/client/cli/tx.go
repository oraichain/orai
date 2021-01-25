package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/oraichain/orai/x/websocket/types"
)

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

	return cmd
}
