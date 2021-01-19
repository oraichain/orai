package cli

import (
	"strings"

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
		GetCmdCreateOracleScript(),
	)
	return txCmd
}

// GetCmdCreateAIDataSource is the CLI command for sending a SetAIDataSource transaction
func GetCmdCreateAIDataSource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-datasource [name] [contract] [description]",
		Short: "Set a new data source into the system",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			// collect transaction fee from the user
			fees := viper.GetString(flags.FlagFees)

			msg := types.NewMsgCreateAIDataSource(args[0], args[1], clientCtx.GetFromAddress(), fees, args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCreateAIDataSource is the CLI command for sending a SetAIDataSource transaction
func GetCmdCreateOracleScript() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-oscript [name] [contract] [description]",
		Short: "Set a new oracle script into the system",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// should find better way
			dSources := strings.Split(args[2], ",")
			tCases := strings.Split(args[3], ",")

			// collect transaction fee from the user
			fees := viper.GetString(flags.FlagFees)

			msg := types.NewMsgCreateOracleScript(args[0], args[1], clientCtx.GetFromAddress(), fees, dSources, tCases)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
