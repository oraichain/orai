package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		Use:   "set-orai [contract-address] [input] [request-fees] [validator-count] [test-only]",
		Short: "Set a new ai oracle request and set result into the system, test only should be either true or false",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valCount, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			testOnly, err := strconv.ParseBool(args[4])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetAIRequestReq(ksuid.New().String(), contractAddr.String(), clientCtx.GetFromAddress().String(), args[2], int64(valCount), []byte(args[1]), testOnly)
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
