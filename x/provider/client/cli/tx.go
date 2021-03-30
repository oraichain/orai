package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/oraichain/orai/x/provider/types"
	"github.com/spf13/cobra"
)

const (
	flagDataSources = "ds"
	flagTestCases   = "tc"
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
		GetCmdSetAIDataSource(),
		GetCmdSetOracleScript(),
		GetCmdSetTestCase(),
		GetCmdEditOracleScript(),
		GetCmdEditAIDataSource(),
		GetCmdEditTestCase(),
	)
	return txCmd
}

// GetCmdSetAIDataSource is the CLI command for sending a SetAIDataSource transaction
func GetCmdSetAIDataSource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-datasource [name] [contract] [description] [fees]",
		Short: "Set a new data source into the system",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAIDataSource(args[0], args[1], clientCtx.GetFromAddress(), args[3], args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdSetOracleScript is the CLI command for sending a SetOracleScript transaction
func GetCmdSetOracleScript() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-oscript [name] [contract] [description] [fees] (--ds [datasource list]) (--tc [testcase list])",
		Short: "Set a new oscript into the system",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}

			// get testcases and datasources
			dSources, err := cmd.Flags().GetStringSlice(flagDataSources)
			if err != nil {
				return err
			}

			tCases, err := cmd.Flags().GetStringSlice(flagTestCases)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateOracleScript(args[0], args[1], clientCtx.GetFromAddress(), args[3], args[2], dSources, tCases)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}

	cmd.Flags().StringSlice(flagDataSources, make([]string, 0), "identifiers of the data sources")
	cmd.Flags().StringSlice(flagTestCases, make([]string, 0), "identifiers of the test cases")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdSetTestCase is the CLI command for sending a SetTestCase transaction
func GetCmdSetTestCase() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-testcase [name] [contract] [description] [fees]",
		Short: "Set a new ai request test case into the system",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateTestCase(args[0], args[1], clientCtx.GetFromAddress(), args[3], args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdEditAIDataSource is the CLI command for sending a EditAIDataSource transaction
func GetCmdEditAIDataSource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-datasource [old-name] [new-name] [contract] [description] [fees]",
		Short: "Edit an existing data source in the system",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditAIDataSource(args[0], args[1], args[2], clientCtx.GetFromAddress(), args[4], args[3])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdEditOracleScript is the CLI command for sending a EditOracleScript transaction
func GetCmdEditOracleScript() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-oscript [old-name] [new-name] [contract] [description] [fees] (--ds [datasource list]) (--tc [testcase list])",
		Short: "Edit an existing oscript in the system",
		Args:  cobra.MinimumNArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// get testcases and datasources
			dSources, err := cmd.Flags().GetStringSlice(flagDataSources)
			if err != nil {
				return err
			}

			tCases, err := cmd.Flags().GetStringSlice(flagTestCases)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditOracleScript(args[0], args[1], args[2], clientCtx.GetFromAddress(), args[4], args[3], dSources, tCases)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice(flagDataSources, make([]string, 0), "identifiers of the data sources")
	cmd.Flags().StringSlice(flagTestCases, make([]string, 0), "identifiers of the test cases")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditTestCase is the CLI command for sending a EditTestCase transaction
func GetCmdEditTestCase() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-testcase [old-name] [new-name] [contract] [description] [fees]",
		Short: "Edit an existing test case in the system",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditTestCase(args[0], args[1], args[2], clientCtx.GetFromAddress(), args[4], args[3])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
