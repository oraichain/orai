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

			return clientCtx.PrintBytes(res.Data)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
