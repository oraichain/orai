package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
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
	)
	return txCmd
}

// GetCmdSetAIOracle is the CLI command for sending a SetAIOracle transaction
func GetCmdSetAIOracle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-orai [contract-address] [input] [request-fees] [validator-count]",
		Short: "Set a new ai oracle request and set result into the system",
		Args:  cobra.ExactArgs(4),
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

			msg := types.NewMsgSetAIOracleReq(ksuid.New().String(), contractAddr, clientCtx.GetFromAddress(), args[2], int64(valCount), []byte(args[1]))
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
