package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/oraichain/orai/x/airequest/types"
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
		GetCmdQueryAIRequest(),
		GetCmdAIRequestIDs(),
	)
	return queryCmd
}

// GetCmdQueryAIRequest query an ai request
func GetCmdQueryAIRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aireq [id]",
		Short: "query an AI request",
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
			res, err := queryClient.QueryAIRequest(
				context.Background(),
				&types.QueryAIRequestReq{
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

// GetCmdAIRequestIDs queries a Queryall ai request ids
func GetCmdAIRequestIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aireqs --page [1] --limit [1]",
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
			res, err := queryClient.QueryAIRequestIDs(
				context.Background(),
				&types.QueryAIRequestIDsReq{
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
