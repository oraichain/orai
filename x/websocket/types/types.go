package types

import (
	"time"
)

// WebSocketConfig is the extra config required for wasm
type WebSocketConfig struct {
	BroadcastTimeout time.Duration
	// MemoryCacheSize in MiB not bytes
	RPCPollInterval time.Duration
	// ContractDebugMode log what contract print
	MaxTry uint64
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWebSocketConfig() WebSocketConfig {
	return WebSocketConfig{
		BroadcastTimeout: time.Minute * 5,
		RPCPollInterval:  time.Second,
		MaxTry:           5,
	}
}
