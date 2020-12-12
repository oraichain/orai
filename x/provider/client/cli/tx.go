package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/oraichain/orai/x/provider/types"
)

const (
	flagDataSources = "ds"
	flagTestCases   = "tc"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	providerTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	providerTxCmd.AddCommand(flags.PostCommands(
		// TODO: Add tx based commands
		// GetCmd<Action>(cdc)
		GetCmdSetOracleScript(cdc),
		GetCmdEditOracleScript(cdc),
		GetCmdCreateAIDataSource(cdc),
		GetCmdEditAIDataSource(cdc),
		//GetCmdSetAIRequest(cdc),
		GetCmdCreateTestCase(cdc),
		GetCmdEditTestCase(cdc),
		// GetCmdAddReporters(cdc),
		// GetCmdRemoveReporter(cdc),
		// GetCmdCreateStrategy(cdc),
	)...)

	return providerTxCmd
}

// Example:
//
// GetCmd<Action> is the CLI command for doing <Action>
// func GetCmd<Action>(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "/* Describe your action cmd */",
// 		Short: "/* Provide a short description on the cmd */",
// 		Args:  cobra.ExactArgs(2), // Does your request require arguments
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			msg := types.NewMsg<Action>(/* Action params */)
// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }

// GetCmdSetOracleScript is the CLI command for sending a SetOracleScript transaction
func GetCmdSetOracleScript(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-oscript [name] [code-path] [description] (--ds [datasource list]) (--tc [testcase list])",
		Short: "Set a new oscript into the system",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			execBytes, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			dSources, err := cmd.Flags().GetStringSlice(flagDataSources)
			if err != nil {
				return err
			}

			tCases, err := cmd.Flags().GetStringSlice(flagTestCases)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateOracleScript(args[0], execBytes, cliCtx.FromAddress, args[2], dSources, tCases)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringSlice(flagDataSources, make([]string, 0), "identifiers of the data sources")
	cmd.Flags().StringSlice(flagTestCases, make([]string, 0), "identifiers of the test cases")

	return cmd
}

// GetCmdEditOracleScript is the CLI command for sending a MsgEditOracleScript transaction
func GetCmdEditOracleScript(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-oscript [old-name] [new-name] [code-path] [description]",
		Short: "Edit an existing oscript in the system",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			execBytes, err := ioutil.ReadFile(args[2])
			if err != nil {
				return err
			}

			dSources, err := cmd.Flags().GetStringSlice(flagDataSources)
			if err != nil {
				return err
			}

			tCases, err := cmd.Flags().GetStringSlice(flagTestCases)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditOracleScript(args[0], args[1], execBytes, cliCtx.FromAddress, args[3], dSources, tCases)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringSlice(flagDataSources, make([]string, 0), "identifiers of the data sources")
	cmd.Flags().StringSlice(flagTestCases, make([]string, 0), "identifiers of the test cases")

	return cmd
}

// GetCmdCreateAIDataSource is the CLI command for sending a SetAIDataSource transaction
func GetCmdCreateAIDataSource(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-datasource [name] [code-path] [description]",
		Short: "Set a new data source into the system",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			execBytes, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			// collect transaction fee from the user
			fees := viper.GetString(flags.FlagFees)

			msg := types.NewMsgCreateAIDataSource(args[0], execBytes, cliCtx.GetFromAddress(), fees, args[2])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdEditAIDataSource is the CLI command for sending a MsgEditDataSource transaction
func GetCmdEditAIDataSource(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "edit-datasource [old-name] [new-name] [code-path] [description]",
		Short: "Edit an existing data source in the system",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			execBytes, err := ioutil.ReadFile(args[2])
			if err != nil {
				return err
			}

			// collect transaction fee from the user
			fees := viper.GetString(flags.FlagFees)

			msg := types.NewMsgEditAIDataSource(args[0], args[1], execBytes, cliCtx.FromAddress, fees, args[3])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// // GetCmdSetAIRequest is the CLI command for sending a SetAIRequest transaction
// func GetCmdSetAIRequest(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "set-aireq [name] [image-path]",
// 		Short: "Set a new ai request and set result into the system",
// 		Args:  cobra.ExactArgs(2),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			imageBytes, err := ioutil.ReadFile(args[1])
// 			if err != nil {
// 				return err
// 			}

// 			// Has to compress the image since it's too large
// 			var buf bytes.Buffer
// 			zw := gzip.NewWriter(&buf)
// 			zw.Write(imageBytes)
// 			zw.Close()
// 			imageCompressed := buf.Bytes()

// 			msg := types.NewMsgSetAIRequest(ksuid.New().String(), args[0], cliCtx.GetFromAddress(), imageCompressed, args[1])
// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }

// GetCmdCreateTestCase is the CLI command for sending a SetTestCase transaction
func GetCmdCreateTestCase(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-testcase [name] [code-path] [description]",
		Short: "Set a new ai request test case into the store",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			execBytes, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			// collect transaction fee from the user
			fees := viper.GetString(flags.FlagFees)

			msg := types.NewMsgCreateTestCase(args[0], execBytes, cliCtx.GetFromAddress(), fees, args[2])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdEditTestCase is the CLI command for sending a MsgEditTestCase transaction
func GetCmdEditTestCase(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "edit-testcase [old-name] [new-name] [code-path] [description]",
		Short: "Edit an existing data source in the system",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			execBytes, err := ioutil.ReadFile(args[2])
			if err != nil {
				return err
			}

			// collect transaction fee from the user
			fees := viper.GetString(flags.FlagFees)

			msg := types.NewMsgEditTestCase(args[0], args[1], execBytes, cliCtx.FromAddress, fees, args[3])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// // GetCmdAddReporters implements the add reporters command handler.
// func GetCmdAddReporters(cdc *codec.Codec) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "add-reporters [reporter1] [reporter2] ...",
// 		Short: "Add agents authorized to submit report transactions.",
// 		Args:  cobra.MinimumNArgs(1),
// 		Long: strings.TrimSpace(
// 			fmt.Sprintf(`Add agents authorized to submit report transactions.
// Example:
// $ %s tx provider add-reporters orai1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun orai1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs --from mykey
// `,
// 				version.ClientName,
// 			),
// 		),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			validator := sdk.ValAddress(cliCtx.GetFromAddress())
// 			msgs := make([]sdk.Msg, len(args))
// 			for i, raw := range args {
// 				reporter, err := sdk.AccAddressFromBech32(raw)
// 				if err != nil {
// 					return err
// 				}
// 				msgs[i] = types.NewMsgAddReporter(
// 					validator,
// 					reporter,
// 					cliCtx.GetFromAddress(),
// 				)
// 				err = msgs[i].ValidateBasic()
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
// 		},
// 	}

// 	return cmd
// }

// // GetCmdRemoveReporter implements the remove reporter command handler.
// func GetCmdRemoveReporter(cdc *codec.Codec) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "remove-reporter [reporter]",
// 		Short: "Remove an agent from the list of authorized reporters.",
// 		Args:  cobra.ExactArgs(1),
// 		Long: strings.TrimSpace(
// 			fmt.Sprintf(`Remove an agent from the list of authorized reporters.
// Example:
// $ %s tx provider remove-reporter band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun --from mykey
// `,
// 				version.ClientName,
// 			),
// 		),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
// 			validator := sdk.ValAddress(cliCtx.GetFromAddress())
// 			reporter, err := sdk.AccAddressFromBech32(args[0])
// 			if err != nil {
// 				return err
// 			}
// 			msg := types.NewMsgRemoveReporter(
// 				validator,
// 				reporter,
// 				cliCtx.GetFromAddress(),
// 			)
// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}
// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}

// 	return cmd
// }

// // GetCmdCreateStrategy is the CLI command for creating a new strategy transaction
// func GetCmdCreateStrategy(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "set-strat [id] [name] [performance-fee] [performance-max] [withdrawal-fee] [withdrawal-max] [gov-addr] [strategist-addr] [flow] ...",
// 		Short: "Set a new strategy into the store",
// 		Args:  cobra.MinimumNArgs(9),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			var flow []string

// 			for i := 8; i < len(args); i++ {
// 				flow = append(flow, args[i])
// 			}

// 			// convert some int values from string to uint64
// 			id, err := strconv.Atoi(args[0])
// 			if err != nil {
// 				return err
// 			}

// 			perFee, err := strconv.Atoi(args[2])
// 			if err != nil {
// 				return err
// 			}

// 			perMax, err := strconv.Atoi(args[3])
// 			if err != nil {
// 				return err
// 			}

// 			wdrawFee, err := strconv.Atoi(args[4])
// 			if err != nil {
// 				return err
// 			}

// 			wdrawMax, err := strconv.Atoi(args[5])
// 			if err != nil {
// 				return err
// 			}

// 			// create a new MsgStrategy
// 			msg := types.NewMsgCreateStrategy(uint64(id), args[1], flow, uint64(perFee), uint64(perMax), uint64(wdrawFee), uint64(wdrawMax), args[6], args[7], cliCtx.GetFromAddress())

// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }
