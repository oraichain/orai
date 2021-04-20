package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/oraichain/orai/x/aioracle/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagPage     = "page"
	flagLimit    = "limit"
	flagValNum   = "val-num"
	flagTestOnly = "test-only"
	LineBreak    = byte('\n')
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// using the same clientRequest as rest
	queryCmd.AddCommand(
		GetCmdQueryAIOracle(),
		GetCmdAIOracleIDs(),
		GetCmdQueryDataSourceContract(),
		//GetCmdQueryTestCaseContract(),
		GetCmdQueryOScriptContract(),
		GetCmdQueryTestCaseEntries(),
		GetCmdQueryDataSourceEntries(),
		GetCmdQueryFullRequest(),
		GetCmdReward(),
		GetCmdQueryMinFees(),
	)
	return queryCmd
}

// GetCmdQueryAIOracle query an ai request
func GetCmdQueryAIOracle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orai [id]",
		Short: "query an AI Oracle request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// validate basic
			if len(args[0]) == 0 {
				return errors.New("request ID can not be empty")
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryAIOracle(
				context.Background(),
				&types.QueryAIOracleReq{
					RequestId: args[0],
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

// GetCmdAIOracleIDs queries a Queryall ai request ids
func GetCmdAIOracleIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orais --page [1] --limit [1]",
		Short: "query all AI request IDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			page, err := cmd.Flags().GetInt64(flagPage)
			if err != nil {
				return err
			}

			limit, err := cmd.Flags().GetInt64(flagLimit)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryAIOracleIDs(
				context.Background(),
				&types.QueryAIOracleIDsReq{
					Page:  page,
					Limit: limit,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Int64(flagPage, types.DefaultQueryPage, "from page")
	cmd.Flags().Int64(flagLimit, types.DefaultQueryLimit, "limit number")
	return cmd
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

// // GetCmdQueryTestCaseContract lists data source code uploaded
// func GetCmdQueryTestCaseContract() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "tcase [address] [entrypoint] [input] [output]",
// 		Short: "query an testcase smart contract",
// 		Args:  cobra.ExactArgs(4),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientQueryContext(cmd)
// 			if err != nil {
// 				return err
// 			}

// 			queryClient := types.NewQueryClient(clientCtx)
// 			contractAddr, err := sdk.AccAddressFromBech32(args[0])
// 			if err != nil {
// 				return err
// 			}

// 			var entrypoint types.EntryPoint
// 			if err := clientCtx.JSONMarshaler.UnmarshalJSON([]byte(args[1]), &entrypoint); err != nil {
// 				return err
// 			}

// 			res, err := queryClient.TestCaseContract(
// 				context.Background(),
// 				&types.QueryTestCaseContract{
// 					Contract: contractAddr,
// 					Request: &types.RequestTestCase{
// 						Tcase:  &entrypoint,
// 						Input:  args[2],
// 						Output: args[3],
// 					},
// 				},
// 			)

// 			if err != nil {
// 				return err
// 			}

// 			return clientCtx.PrintBytes(append(res.Data, LineBreak))
// 		},
// 	}
// 	flags.AddQueryFlagsToCmd(cmd)
// 	return cmd
// }

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

// GetCmdQueryFullRequest queries full information about an AI request
func GetCmdQueryFullRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fullreq [id]",
		Short: "query a full ai request using its request ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			id := args[0]

			queryClient := types.NewQueryClient(cliCtx)
			res, err := queryClient.QueryFullRequest(
				context.Background(),
				&types.QueryFullOracleReq{
					RequestId: id,
				},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdReward queries full information about a reward at a specific block height
func GetCmdReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward [block-height]",
		Short: "query a reward information given a block height",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			blockHeight := args[0]

			queryClient := types.NewQueryClient(cliCtx)
			res, err := queryClient.QueryReward(
				context.Background(),
				&types.QueryRewardReq{
					BlockHeight: blockHeight,
				},
			)
			if err != nil {
				return err
			}

			return cliCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryMinFees queries a the minimum fees of an ai oracle request
func GetCmdQueryMinFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "minfees [contract-addr] [1] [true/false]",
		Short: "query the min fees",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			valNum, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			testOnly, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryMinFees(
				context.Background(),
				&types.MinFeesReq{
					ContractAddr: args[0],
					ValNum:       int64(valNum),
					TestOnly:     testOnly,
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
