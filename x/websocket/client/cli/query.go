package cli

import (
	"context"
	"fmt"

	"github.com/oraichain/orai/x/websocket/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const LineBreak = byte('\n')

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group provider queries under a subcommand
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryOracleInfo(),
		GetCmdQueryDataSourceContract(),
		GetCmdQueryTestCaseContract(),
		GetCmdQueryOScriptContract(),
	)

	return queryCmd
}

// GetCmdQueryOracleInfo lists data source code uploaded
func GetCmdQueryOracleInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle [contract] [url]",
		Short: "query an oracle smart contract",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			contract := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			contractAddr, err := sdk.AccAddressFromBech32(contract)
			if err != nil {
				return err
			}

			res, err := queryClient.OracleContract(
				context.Background(),
				&types.QueryOracleContract{
					Contract: contractAddr,
					Request: &types.Request{
						Fetch: &types.Fetch{
							Url: args[1],
						},
					},
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintBytes(append(res.Data, LineBreak))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryDataSourceContract lists data source code uploaded
func GetCmdQueryDataSourceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dsource [name] [input]",
		Short: "query an datasource smart contract",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			input := ""
			if len(args) > 1 {
				input = args[1]
			}
			res, err := queryClient.DataSourceContract(
				context.Background(),
				&types.QueryDataSourceContract{
					Name: args[0],
					Request: &types.RequestDataSource{
						Input: input,
					},
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintBytes(append(res.Data, LineBreak))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTestCaseContract lists data source code uploaded
func GetCmdQueryTestCaseContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tcase [name] [dsource] [input] [output]",
		Short: "query an testcase smart contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TestCaseContract(
				context.Background(),
				&types.QueryTestCaseContract{
					Name:           args[0],
					DataSourceName: args[1],
					Request: &types.RequestTestCase{
						Input:  args[2],
						Output: args[3],
					},
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintBytes(append(res.Data, LineBreak))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryOScriptContract lists data source code uploaded
func GetCmdQueryOScriptContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oscript [name] [results...]",
		Short: "query an oscript smart contract",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			query := &types.QueryOracleScriptContract{
				Name: args[0],
				Request: &types.RequestOracleScript{
					Results: args[1:],
				},
			}
			fmt.Printf("query :%v\n", query)
			res, err := queryClient.OracleScriptContract(context.Background(), query)

			if err != nil {
				return err
			}

			return clientCtx.PrintBytes(append(res.Data, LineBreak))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
