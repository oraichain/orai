package websocket

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	keyring "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// TxQuery ...
	TxQuery = "tm.event = 'Tx'"
	// EventChannelCapacity is a buffer size of channel between node and this program
	EventChannelCapacity = 5000
)

func runImpl(cdc codec.Marshaler, c *Context, l *Logger) error {
	l.Info(":rocket: Starting WebSocket subscriber")

	ctx, cxl := context.WithTimeout(context.Background(), 100*time.Second)
	defer cxl()

	// start listening to the events from the 26657 port after creating the container successfully
	err := c.client.Start()
	if err != nil {
		return err
	}
	eventChan, err := c.client.Subscribe(ctx, "", TxQuery, EventChannelCapacity)
	if err != nil {
		return err
	}

	for {
		select {
		case ev := <-eventChan:
			l.Info("%v\n", ev.Data.(tmtypes.EventDataTx).TxResult)
			go handleTransaction(c, l, ev.Data.(tmtypes.EventDataTx).TxResult)
		}
	}
}

func RegisterWebSocketFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(flagValidator, "", "validator address")
	// cmd.Flags().String(flagExecutor, "", "executor name and url for executing the data source script")
	cmd.Flags().String(flags.FlagGasPrices, "", "gas prices for report transaction")
	cmd.Flags().String(flagLogLevel, "info", "set the logger level")
	cmd.Flags().String(flagBroadcastTimeout, "5m", "The time that the websocket will wait for tx commit")
	cmd.Flags().String(flagRPCPollInterval, "1s", "The duration of rpc poll interval")
	cmd.Flags().Uint64(flagMaxTry, 5, "The maximum number of tries to submit a report transaction")
	viper.BindPFlag(flagValidator, cmd.Flags().Lookup(flagValidator))
	viper.BindPFlag(flags.FlagGasPrices, cmd.Flags().Lookup(flags.FlagGasPrices))
	viper.BindPFlag(flagLogLevel, cmd.Flags().Lookup(flagLogLevel))
	//viper.BindPFlag(flagExecutor, cmd.Flags().Lookup(flagExecutor))
	viper.BindPFlag(flagBroadcastTimeout, cmd.Flags().Lookup(flagBroadcastTimeout))
	viper.BindPFlag(flagRPCPollInterval, cmd.Flags().Lookup(flagRPCPollInterval))
	viper.BindPFlag(flagMaxTry, cmd.Flags().Lookup(flagMaxTry))
	return cmd
}

func runCmd(cdc codec.Marshaler, c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Short:   "Run the oracle process",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.ChainID == "" {
				return errors.New("Chain ID must not be empty")
			}

			gasAdj, err := cmd.Flags().GetFloat64(flags.FlagGasAdjustment) // 1
			if err != nil {
				gasAdj, err = strconv.ParseFloat(cfg.GasAdjustment, 64)
				if err != nil || gasAdj == float64(0) {
					gasAdj = flags.DefaultGasAdjustment
				}
			}
			c.gasAdj = gasAdj
			gas, err := cmd.Flags().GetUint64("gas") // 200000
			if err != nil {
				gasInt, err := strconv.Atoi(cfg.Gas)
				if err != nil || gas == uint64(0) {
					gas = flags.DefaultGasLimit
				} else {
					gas = uint64(gasInt)
				}
			}
			c.gas = gas
			feesStr, err := cmd.Flags().GetString(flags.FlagFees) // 5000orai
			if err != nil {
				feesStr = defaultFees
				c.fees, _ = sdk.ParseCoinsNormalized(defaultFees)
			} else {
				fees, err := sdk.ParseCoinsNormalized(feesStr)
				if err != nil {
					fees, _ = sdk.ParseCoinsNormalized(defaultFees)
				}
				c.fees = fees
			}

			// other params
			keys, err := keybase.List()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				return errors.New("No key available")
			}
			c.keys = make(chan keyring.Info, len(keys))
			for _, key := range keys {
				c.keys <- key
			}
			c.validator, err = sdk.ValAddressFromBech32(cfg.Validator)
			if err != nil {
				return err
			}
			err = sdk.VerifyAddressFormat(c.validator)
			if err != nil {
				return err
			}
			c.gasPrices, err = sdk.ParseDecCoins(cfg.GasPrices)
			if err != nil {
				return err
			}

			l.Info(":star: Creating the daemon listening to node: %s", cfg.NodeURI)
			c.client, err = httpclient.New(cfg.NodeURI, "/websocket")
			if err != nil {
				return err
			}
			// c.fileCache = filecache.New(filepath.Join(viper.GetString(flags.FlagHome), "files"))
			c.broadcastTimeout, err = time.ParseDuration(cfg.BroadcastTimeout)
			if err != nil {
				return err
			}
			c.maxTry = cfg.MaxTry
			c.rpcPollInterval, err = time.ParseDuration(cfg.RPCPollInterval)
			if err != nil {
				return err
			}
			return runImpl(cdc, c, l)
		},
	}

	return RegisterWebSocketFlags(cmd)

}
