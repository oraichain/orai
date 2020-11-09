package cli

import (
	"fmt"

	"github.com/oraichain/orai/x/airequest/types"
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
			GetCmdQueryAIRequest(queryRoute, cdc),
			GetCmdAIRequestIDs(queryRoute, cdc),
		)...,
	)

	return providerQueryCmd
}

// TODO: Add Query Commands

// GetCmdQueryAIRequest queries information about an AI request
func GetCmdQueryAIRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "aireq [id]",
		Short: "query an ai request using its request ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := []byte(args[0])

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireq/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Printf("could not query request - %s \n", args[0])
				return nil
			}

			var out types.QueryResAIRequest
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdAIRequestIDs queries a list of all request IDs
func GetCmdAIRequestIDs(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "aireqs",
		Short: "query all AI request IDs",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/aireqs", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get request IDs\n")
				return nil
			}

			var out types.QueryResAIRequestIDs
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
