package app

import (
	"encoding/json"
	"time"

	wasm "github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidence "github.com/cosmos/cosmos-sdk/x/evidence/types"
	genutil "github.com/cosmos/cosmos-sdk/x/genutil/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/core/types"
	mint "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	aiRequest "github.com/oraichain/orai/x/airequest/types"
	airesult "github.com/oraichain/orai/x/airesult/types"
	provider "github.com/oraichain/orai/x/provider/types"
	websocket "github.com/oraichain/orai/x/websocket/types"
)

// GenesisState default state for the application
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONMarshaler) GenesisState {
	genesisState := make(GenesisState)
	// Get default genesis states of the modules we are to override.
	authGenesis := auth.DefaultGenesisState()
	stakingGenesis := staking.DefaultGenesisState()
	distrGenesis := distr.DefaultGenesisState()
	mintGenesis := mint.DefaultGenesisState()
	govGenesis := gov.DefaultGenesisState()
	crisisGenesis := crisis.DefaultGenesisState()
	slashingGenesis := slashing.DefaultGenesisState()
	// Override the genesis parameters.
	authGenesis.Params.TxSizeCostPerByte = 5
	stakingGenesis.Params.BondDenom = provider.Denom
	stakingGenesis.Params.HistoricalEntries = 1000

	// TODO: testnet figures only
	stakingGenesis.Params.UnbondingTime = time.Hour * 2
	stakingGenesis.Params.MaxValidators = 100
	// maximum bonded validators
	distrGenesis.Params.BaseProposerReward = sdk.NewDecWithPrec(1, 2)  // 1%
	distrGenesis.Params.BonusProposerReward = sdk.NewDecWithPrec(4, 2) // 4%
	mintGenesis.Params.BlocksPerYear = 6311200                         // target 5-second block time
	mintGenesis.Params.MintDenom = provider.Denom
	govGenesis.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(provider.Denom, sdk.TokensFromConsensusPower(10)))
	// TODO: testing
	govGenesis.VotingParams.VotingPeriod = time.Minute * 2 // test for 10 mins voting period
	crisisGenesis.ConstantFee = sdk.NewCoin(provider.Denom, sdk.TokensFromConsensusPower(10))
	slashingGenesis.Params.SignedBlocksWindow = 30000                         // approximately 1 day
	slashingGenesis.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2)      // 5%
	slashingGenesis.Params.DowntimeJailDuration = 60 * 10 * time.Second       // 10 minutes
	slashingGenesis.Params.SlashFractionDoubleSign = sdk.NewDecWithPrec(5, 2) // 5%
	slashingGenesis.Params.SlashFractionDowntime = sdk.NewDecWithPrec(1, 4)   // 0.01%
	// Add your modules here for the genesis states

	genesisState[genutil.ModuleName] = cdc.MustMarshalJSON(genutil.DefaultGenesisState())
	genesisState[auth.ModuleName] = cdc.MustMarshalJSON(authGenesis)
	genesisState[bank.ModuleName] = cdc.MustMarshalJSON(bank.DefaultGenesisState())
	genesisState[staking.ModuleName] = cdc.MustMarshalJSON(stakingGenesis)
	genesisState[mint.ModuleName] = cdc.MustMarshalJSON(mintGenesis)
	genesisState[distr.ModuleName] = cdc.MustMarshalJSON(distrGenesis)
	genesisState[gov.ModuleName] = cdc.MustMarshalJSON(govGenesis)
	genesisState[crisis.ModuleName] = cdc.MustMarshalJSON(crisisGenesis)
	genesisState[capabilitytypes.ModuleName] = cdc.MustMarshalJSON(capabilitytypes.DefaultGenesis())
	genesisState[ibchost.ModuleName] = cdc.MustMarshalJSON(ibc.DefaultGenesisState())
	genesisState[ibctransfertypes.ModuleName] = cdc.MustMarshalJSON(ibctransfertypes.DefaultGenesisState())
	genesisState[slashing.ModuleName] = cdc.MustMarshalJSON(slashingGenesis)
	genesisState[evidence.ModuleName] = cdc.MustMarshalJSON(evidence.DefaultGenesisState())
	genesisState[provider.ModuleName] = cdc.MustMarshalJSON(provider.DefaultGenesisState())
	genesisState[aiRequest.ModuleName] = cdc.MustMarshalJSON(aiRequest.DefaultGenesisState())
	genesisState[websocket.ModuleName] = cdc.MustMarshalJSON(websocket.DefaultGenesisState())
	genesisState[airesult.ModuleName] = cdc.MustMarshalJSON(airesult.DefaultGenesisState())
	genesisState[wasm.ModuleName] = cdc.MustMarshalJSON(&wasm.GenesisState{
		Params: wasm.DefaultParams(),
	})

	return genesisState
}
