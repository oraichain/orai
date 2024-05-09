package config

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	Bech32Prefix      = "orai"
	EvmDenom          = Bech32Prefix
	EvmInitialBaseFee = sdk.NewIntFromUint64(1) // 1uorai, still relatively big since initial base fee of eth is 100gwei = 10^9 wei -> smaller than 
)
