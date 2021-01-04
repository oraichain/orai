package app

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	aiRequest "github.com/oraichain/orai/x/airequest"
	aiResult "github.com/oraichain/orai/x/airesult"
	"github.com/oraichain/orai/x/provider"
	"github.com/oraichain/orai/x/provider/types"
	webSocket "github.com/oraichain/orai/x/websocket"
)

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	cdc := MakeCodec()
	denom := types.Denom
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
		genutil.ModuleName:   genutil.AppModuleBasic{}.DefaultGenesis(),
		auth.ModuleName:      cdc.MustMarshalJSON(authGenesis),
		bank.ModuleName:      bank.AppModuleBasic{}.DefaultGenesis(),
		supply.ModuleName:    supply.AppModuleBasic{}.DefaultGenesis(),
		staking.ModuleName:   cdc.MustMarshalJSON(stakingGenesis),
		mint.ModuleName:      cdc.MustMarshalJSON(mintGenesis),
		distr.ModuleName:     cdc.MustMarshalJSON(distrGenesis),
		gov.ModuleName:       cdc.MustMarshalJSON(govGenesis),
		crisis.ModuleName:    cdc.MustMarshalJSON(crisisGenesis),
		slashing.ModuleName:  cdc.MustMarshalJSON(slashingGenesis),
		evidence.ModuleName:  evidence.AppModuleBasic{}.DefaultGenesis(),
		provider.ModuleName:  provider.AppModuleBasic{}.DefaultGenesis(),
		aiRequest.ModuleName: aiRequest.AppModuleBasic{}.DefaultGenesis(),
		webSocket.ModuleName: webSocket.AppModuleBasic{}.DefaultGenesis(),
		aiResult.ModuleName:  aiResult.AppModuleBasic{}.DefaultGenesis(),
	}
}
