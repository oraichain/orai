package cli

import (
	"context"
	"fmt"

	"github.com/oraichain/orai/x/provider/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetAIDataSource(),
		GetOracleScript(),
	)
	return queryCmd
}

// GetAIDataSource lists data source code uploaded
func GetAIDataSource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datasource [name]",
		Short: "query an AI data source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.DataSourceInfo(
				context.Background(),
				&types.DataSourceInfoReq{
					Name: name,
				},
			)
			if err != nil {
				return err
			}
			if len(res.Name) == 0 {
				return fmt.Errorf("data source not found")
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetOracleScript queries information about a oScript
func GetOracleScript() *cobra.Command {
	return &cobra.Command{
		Use:   "oscript [name]",
		Short: "query oscript",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.OracleScriptInfo(
				context.Background(),
				&types.OracleScriptInfoReq{
					Name: name,
				},
			)
			if err != nil {
				return err
			}
			if len(res.Name) == 0 {
				return fmt.Errorf("oscript not found")
			}

			return clientCtx.PrintProto(res)
		},
	}
}
