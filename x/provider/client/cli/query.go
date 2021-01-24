package cli

import (
	"context"
	"fmt"

	"github.com/oraichain/orai/x/provider/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

const (
	flagPage   = "page"
	flagLimit  = "limit"
	flagValNum = "val_num"
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
		GetCmdQueryDataSource(),
		GetCmdQueryDataSources(),
		GetCmdQueryOracleScript(),
		GetCmdQueryOracleScripts(),
		GetCmdQueryTestCase(),
		GetCmdQueryTestCases(),
		GetCmdQueryMinFees(),
	)
	return queryCmd
}

// GetCmdQueryDataSource lists data source code uploaded
func GetCmdQueryDataSource() *cobra.Command {
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

// GetCmdQueryDataSources queries a Queryall data source names
func GetCmdQueryDataSources() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dsources [name] --page [1] --limit [5]",
		Short: "query all AI data sources",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			name := ""
			if len(args) > 0 {
				name = args[0]
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
			res, err := queryClient.ListDataSources(
				context.Background(),
				&types.ListDataSourcesReq{
					Name:  name,
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

// GetCmdQueryOracleScript queries information about a oScript
func GetCmdQueryOracleScript() *cobra.Command {
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

// GetCmdQueryOracleScripts queries a list of all oscript names
func GetCmdQueryOracleScripts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oscripts",
		Short: "query all oscripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			name := ""
			if len(args) > 0 {
				name = args[0]
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
			res, err := queryClient.ListOracleScripts(
				context.Background(),
				&types.ListOracleScriptsReq{
					Name:  name,
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

// GetCmdQueryTestCase queries information about an AI request test case
func GetCmdQueryTestCase() *cobra.Command {
	return &cobra.Command{
		Use:   "testcase [name]",
		Short: "query an ai request test case using name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TestCaseInfo(
				context.Background(),
				&types.TestCaseInfoReq{
					Name: name,
				},
			)
			if err != nil {
				return err
			}
			if len(res.Name) == 0 {
				return fmt.Errorf("testcase not found")
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetCmdQueryTestCases queries a list of all test case names
func GetCmdQueryTestCases() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tcases",
		Short: "query all AI request test cases",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			name := ""
			if len(args) > 0 {
				name = args[0]
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
			res, err := queryClient.ListTestCases(
				context.Background(),
				&types.ListTestCasesReq{
					Name:  name,
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

// GetCmdQueryMinFees queries a list of all test case names
func GetCmdQueryMinFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "minfees [name] --val_num [1]",
		Short: "query the min fees",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			valNum, err := cmd.Flags().GetInt64(flagValNum)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryMinFees(
				context.Background(),
				&types.MinFeesReq{
					OracleScriptName: args[0],
					ValNum:           valNum,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Int64(flagValNum, types.DefaultValNum, "val num")
	return cmd
}
