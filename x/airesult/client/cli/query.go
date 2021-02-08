package cli

import (
	"context"
	"fmt"

	"github.com/oraichain/orai/x/airesult/types"
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

	// using the same clientRequest as rest
	queryCmd.AddCommand(
		GetCmdQueryFullRequest(),
		GetCmdReward(),
	)
	return queryCmd
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
				&types.QueryFullRequestReq{
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
