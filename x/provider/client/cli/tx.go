package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/oraichain/orai/x/provider/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "provider transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdCreateAIDataSource(),
	)
	return txCmd
}

// GetCmdCreateAIDataSource is the CLI command for sending a SetAIDataSource transaction
func GetCmdCreateAIDataSource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-datasource [name] [code-path] [description]",
		Short: "Set a new data source into the system",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			execBytes, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			// collect transaction fee from the user
			fees := viper.GetString(flags.FlagFees)

			msg := types.NewMsgCreateAIDataSource(args[0], execBytes, cliCtx.GetFromAddress(), fees, args[2])
			// err = msg.ValidateBasic()
			// if err != nil {
			// 	return err
			// }

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
