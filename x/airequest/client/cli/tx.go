package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "airequest transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdSetAIRequest(),
	)
	return txCmd
}

// GetCmdSetAIRequest is the CLI command for sending a SetAIRequest transaction
func GetCmdSetAIRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-aireq [oracle-script-name] [input] [expected-output] [request-fees] [validator-count]",
		Short: "Set a new ai request and set result into the system",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valCount, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetAIRequest(ksuid.New().String(), args[0], clientCtx.GetFromAddress(), args[3], int64(valCount), []byte(args[1]), []byte(args[2]))
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
