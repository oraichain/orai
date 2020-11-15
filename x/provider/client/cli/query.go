package cli

import (
	"fmt"

	"github.com/oraichain/orai/x/provider/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group provider queries under a subcommand
	providerQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	providerQueryCmd.AddCommand(
		flags.GetCommands(
			// TODO: Add query Cmds
			GetCmdQueryOracleScript(queryRoute, cdc),
			GetCmdQueryDataSource(queryRoute, cdc),
			GetCmdOracleScriptNames(queryRoute, cdc),
			GetCmdDataSourceNames(queryRoute, cdc),
			// GetCmdQueryAIRequest(queryRoute, cdc),
			// GetCmdAIRequestIDs(queryRoute, cdc),
			GetCmdQueryTestCase(queryRoute, cdc),
			GetCmdTestCaseNames(queryRoute, cdc),
			// GetCmdQueryFullRequest(queryRoute, cdc),
		)...,
	)

	return providerQueryCmd
}

// TODO: Add Query Commands

// GetCmdQueryOracleScript queries information about a oScript
func GetCmdQueryOracleScript(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oscript [name]",
		Short: "query oscript",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/oscript/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not query oScript - %s \n", err.Error())
				return nil
			}

			var out types.QueryResOracleScript
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdOracleScriptNames queries a list of all oscript names
func GetCmdOracleScriptNames(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "onames",
		Short: "query all oscript names",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			fmt.Printf("address string: %s\n", cliCtx.ChainID)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/onames", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get query names\n")
				return nil
			}

			var out types.QueryResOracleScriptNames
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryDataSource queries information about an AIDataSource
func GetCmdQueryDataSource(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "datasource [name]",
		Short: "query datasource",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/datasource/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not query data source - %s \n", name)
				return nil
			}

			var out types.QueryResAIDataSource
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdDataSourceNames queries a list of all data source names
func GetCmdDataSourceNames(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dnames",
		Short: "query all data source names",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/dnames", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get query names\n")
				return nil
			}

			var out types.QueryResAIDataSourceNames
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// // GetCmdQueryAIRequest queries information about an AI request
// func GetCmdQueryAIRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "aireq [id]",
// 		Short: "query an ai request using its request ID",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			id := []byte(args[0])

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireq/%s", queryRoute, id), nil)
// 			if err != nil {
// 				fmt.Printf("could not query request - %s \n", args[0])
// 				return nil
// 			}

// 			var out types.QueryResAIRequest
// 			cdc.MustUnmarshalJSON(res, &out)
// 			return cliCtx.PrintOutput(out)
// 		},
// 	}
// }

// // GetCmdQueryFullRequest queries full information about an AI request
// func GetCmdQueryFullRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "fullreq [id]",
// 		Short: "query a full ai request using its request ID",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			id := []byte(args[0])

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/fullreq/%s", queryRoute, id), nil)
// 			if err != nil {
// 				fmt.Printf("could not query request - %s \n", args[0])
// 				return nil
// 			}

// 			var out types.QueryResFullRequest
// 			cdc.MustUnmarshalJSON(res, &out)
// 			return cliCtx.PrintOutput(out)
// 		},
// 	}
// }

// // GetCmdAIRequestIDs queries a list of all request IDs
// func GetCmdAIRequestIDs(queryRoute string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "aireqs",
// 		Short: "query all AI request IDs",
// 		// Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireqs", queryRoute), nil)
// 			if err != nil {
// 				fmt.Printf("could not get request IDs\n")
// 				return nil
// 			}

// 			var out types.QueryResAIRequestIDs
// 			cdc.MustUnmarshalJSON(res, &out)
// 			return cliCtx.PrintOutput(out)
// 		},
// 	}
// }

// GetCmdQueryTestCase queries information about an AI request test case
func GetCmdQueryTestCase(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "testcase [name]",
		Short: "query an ai request test case using name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := []byte(args[0])

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/testcase/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not query request test case - %s \n", args[0])
				return nil
			}

			var out types.QueryResTestCase
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdTestCaseNames queries a list of all test case names
func GetCmdTestCaseNames(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tcnames",
		Short: "query all AI request test case names",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tcnames", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get test case IDs\n")
				return nil
			}

			var out types.QueryResTestCaseNames
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
