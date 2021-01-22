package app

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidence "github.com/cosmos/cosmos-sdk/x/evidence/types"
	genutil "github.com/cosmos/cosmos-sdk/x/genutil/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	mint "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	aiRequest "github.com/oraichain/orai/x/airequest/types"
	provider "github.com/oraichain/orai/x/provider/types"
)

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	cdc := MakeEncodingConfig()
	denom := "orai"
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
	stakingGenesis.Params.BondDenom = denom
	stakingGenesis.Params.HistoricalEntries = 1000
	// maximum bonded validators
	distrGenesis.Params.BaseProposerReward = sdk.NewDecWithPrec(10, 2)  // 5%
	distrGenesis.Params.BonusProposerReward = sdk.NewDecWithPrec(10, 2) // 12%
	mintGenesis.Params.BlocksPerYear = 6311200                          // target 5-second block time
	mintGenesis.Params.MintDenom = denom
	govGenesis.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(denom, sdk.TokensFromConsensusPower(10)))
	crisisGenesis.ConstantFee = sdk.NewCoin(denom, sdk.TokensFromConsensusPower(10))
	slashingGenesis.Params.SignedBlocksWindow = 30000                         // approximately 1 day
	slashingGenesis.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2)      // 5%
	slashingGenesis.Params.DowntimeJailDuration = 60 * 10 * time.Second       // 10 minutes
	slashingGenesis.Params.SlashFractionDoubleSign = sdk.NewDecWithPrec(5, 2) // 5%
	slashingGenesis.Params.SlashFractionDowntime = sdk.NewDecWithPrec(1, 4)   // 0.01%
	// Add your modules here for the genesis states
	return GenesisState{
		genutil.ModuleName:   cdc.Marshaler.MustMarshalJSON(genutil.DefaultGenesisState()),
		auth.ModuleName:      cdc.Marshaler.MustMarshalJSON(authGenesis),
		bank.ModuleName:      cdc.Marshaler.MustMarshalJSON(bank.DefaultGenesisState()),
		staking.ModuleName:   cdc.Marshaler.MustMarshalJSON(stakingGenesis),
		mint.ModuleName:      cdc.Marshaler.MustMarshalJSON(mintGenesis),
		distr.ModuleName:     cdc.Marshaler.MustMarshalJSON(distrGenesis),
		gov.ModuleName:       cdc.Marshaler.MustMarshalJSON(govGenesis),
		crisis.ModuleName:    cdc.Marshaler.MustMarshalJSON(crisisGenesis),
		slashing.ModuleName:  cdc.Marshaler.MustMarshalJSON(slashingGenesis),
		evidence.ModuleName:  cdc.Marshaler.MustMarshalJSON(evidence.DefaultGenesisState()),
		provider.ModuleName:  cdc.Marshaler.MustMarshalJSON(provider.DefaultGenesisState()),
		aiRequest.ModuleName: cdc.Marshaler.MustMarshalJSON(aiRequest.DefaultGenesisState()),
		// webSocket.ModuleName: webSocket.AppModuleBasic{}.DefaultGenesis(),
		// aiResult.ModuleName:  aiResult.AppModuleBasic{}.DefaultGenesis(),
	}
}
