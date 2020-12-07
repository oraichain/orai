package cli

import (
	"fmt"

	"github.com/oraichain/orai/x/airesult/types"
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
			GetCmdQueryFullRequest(queryRoute, cdc),
			GetCmdReward(queryRoute, cdc),
		)...,
	)

	return providerQueryCmd
}

// TODO: Add Query Commands

// GetCmdQueryFullRequest queries full information about an AI request
func GetCmdQueryFullRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "fullreq [id]",
		Short: "query a full ai request using its request ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := []byte(args[0])

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/fullreq/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Printf("could not query request - %s \n", args[0])
				return nil
			}

			var out types.QueryResFullRequest
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdReward queries full information about a reward at a specific block height
func GetCmdReward(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reward [block-height]",
		Short: "query a reward information given a block height",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			blockHeight := []byte(args[0])

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/reward/%s", queryRoute, blockHeight), nil)
			if err != nil {
				fmt.Printf("could not query request - %s \n", args[0])
				return nil
			}

			var out types.QueryResReward
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
