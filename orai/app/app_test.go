package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"
)

var emptyWasmOpts []wasm.Option = nil

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

func TestWasmdExport(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewOraichainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, MakeEncodingConfig(), wasm.EnableAllProposals, EmptyAppOptions{}, emptyWasmOpts, EvmOptions{})

	genesisState := NewDefaultGenesisState(gapp.appCodec)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	gapp.InitChain(
		abci.RequestInitChain{
			ChainId:       appName,
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	gapp.Commit()

	ctx := gapp.NewContext(true, tmproto.Header{Height: gapp.LastBlockHeight()})
	bz, _ := sdk.GetFromBech32("orai1ur2vsjrjarygawpdwtqteaazfchvw4fg6uql76", "orai")
	creator := sdk.AccAddress(bz)
	wasmCode, _ := os.ReadFile("./bytecode/echo.wasm")
	codeId, _, _ := gapp.ContractKeeper.Create(ctx, creator, wasmCode, nil)
	println("codeid", codeId)

}

// ensure that blocked addresses are properly set in bank keeper
func TestBlockedAddrs(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewOraichainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, MakeEncodingConfig(), wasm.EnableAllProposals, EmptyAppOptions{}, emptyWasmOpts, EvmOptions{})

	for acc := range maccPerms {
		require.Equal(t, !allowedReceivingModAcc[acc], gapp.bankKeeper.BlockedAddr(gapp.accountKeeper.GetModuleAddress(acc)))
	}
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

func TestGetEnabledProposals(t *testing.T) {
	cases := map[string]struct {
		proposalsEnabled string
		specificEnabled  string
		expected         []wasm.ProposalType
	}{
		"all disabled": {
			proposalsEnabled: "false",
			expected:         wasm.DisableAllProposals,
		},
		"all enabled": {
			proposalsEnabled: "true",
			expected:         wasm.EnableAllProposals,
		},
		"some enabled": {
			proposalsEnabled: "okay",
			specificEnabled:  "StoreCode,InstantiateContract",
			expected:         []wasm.ProposalType{wasm.ProposalTypeStoreCode, wasm.ProposalTypeInstantiateContract},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ProposalsEnabled = tc.proposalsEnabled
			EnableSpecificProposals = tc.specificEnabled
			proposals := GetEnabledProposals()
			assert.Equal(t, tc.expected, proposals)
		})
	}
}
