package websocket

import (
	keys "github.com/cosmos/cosmos-sdk/crypto/keyring"
)

const (
	flagValidator        = "validator"
	flagLogLevel         = "log-level"
	flagExecutor         = "executor"
	flagBroadcastTimeout = "broadcast-timeout"
	flagRPCPollInterval  = "rpc-poll-interval"
	flagMaxTry           = "max-try"
	// Bech32MainPrefix is the prefix of different addresses
	Bech32MainPrefix = "orai"

	defaultFees = "5000orai"
)

// Config data structure for the websocket daemon.
type Config struct {
	ChainID          string `mapstructure:"chain-id"`          // ChainID of the target chain
	NodeURI          string `mapstructure:"node"`              // Remote RPC URI of Oraichain node to connect to
	Validator        string `mapstructure:"validator"`         // The validator address that I'm responsible for
	GasPrices        string `mapstructure:"gas-prices"`        // Gas prices of the transaction
	GasAdjustment    string `mapstructure:"gas-adjustment"`    // Gas adjustment from the gas
	Gas              string `mapstructure:"gas"`               // gas used for the transaction
	LogLevel         string `mapstructure:"log-level"`         // Log level of the logger
	Executor         string `mapstructure:"executor"`          // Executor name and URL (example: "Executor name:URL")
	BroadcastTimeout string `mapstructure:"broadcast-timeout"` // The time that the websocket will wait for tx commit
	RPCPollInterval  string `mapstructure:"rpc-poll-interval"` // The duration of rpc poll interval
	MaxTry           uint64 `mapstructure:"max-try"`           // The maximum number of tries to submit a report transaction
}

// Global instances.
var (
	cfg     Config
	keybase keys.Keybase
)

// Main is the entry point of the validator websocket
func Main() {

}
