package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
)

// WebSocketConfig is the extra config required for wasm
type WebSocketConfig struct {
	FromValidator    string
	BroadcastTimeout time.Duration
	RPCPollInterval  time.Duration
	MaxTry           uint64
	Txf              tx.Factory
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWebSocketConfig() *WebSocketConfig {
	return &WebSocketConfig{
		BroadcastTimeout: time.Minute * 5,
		RPCPollInterval:  time.Second,
		MaxTry:           5,
		FromValidator:    "",
	}
}
