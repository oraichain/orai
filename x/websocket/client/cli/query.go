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
		GetCmdQueryDataSourceContract(),
		GetCmdQueryTestCaseContract(),
		GetCmdQueryOScriptContract(),
		GetCmdQueryTestCaseEntries(),
		GetCmdQueryDataSourceEntries(),
	)

	return queryCmd
}

// GetCmdQueryDataSourceContract lists data source code uploaded
func GetCmdQueryDataSourceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dsource [address] [entrypoint] [input]",
		Short: "query an datasource smart contract",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var entrypoint types.EntryPoint
			if err := clientCtx.JSONMarshaler.UnmarshalJSON([]byte(args[1]), &entrypoint); err != nil {
				return err
			}

			input := ""
			if len(args) > 2 {
				input = args[2]
			}

			res, err := queryClient.DataSourceContract(
				context.Background(),
				&types.QueryDataSourceContract{
					Contract: contractAddr,
					Request: &types.RequestDataSource{
						Dsource: &entrypoint,
						Input:   input,
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
		Use:   "tcase [address] [entrypoint] [input] [output]",
		Short: "query an testcase smart contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var entrypoint types.EntryPoint
			if err := clientCtx.JSONMarshaler.UnmarshalJSON([]byte(args[1]), &entrypoint); err != nil {
				return err
			}

			res, err := queryClient.TestCaseContract(
				context.Background(),
				&types.QueryTestCaseContract{
					Contract: contractAddr,
					Request: &types.RequestTestCase{
						Tcase:  &entrypoint,
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
		Use:   "oscript [address] [results...]",
		Short: "query an oscript smart contract",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			query := &types.QueryOracleScriptContract{
				Contract: contractAddr,
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

// GetCmdQueryDataSourceEntries lists data source entries
func GetCmdQueryDataSourceEntries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get_dsources [address]",
		Short: "query data source entries from smart contract",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.DataSourceEntries(
				context.Background(),
				&types.QueryDataSourceEntriesContract{
					Contract: contractAddr,
					Request:  &types.EmptyParams{},
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTestCaseEntries lists data source entries
func GetCmdQueryTestCaseEntries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get_tcases [address]",
		Short: "query test case entries from smart contract",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			contractAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.TestCaseEntries(
				context.Background(),
				&types.QueryTestCaseEntriesContract{
					Contract: contractAddr,
					Request:  &types.EmptyParams{},
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
