package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
)

func TestExport(t *testing.T) {
	db := db.NewMemDB()
	app := NewOraichainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, uint(0), flags.FlagHome)
	setGenesis(app)

	// Making a new app object with the db, so that initchain hasn't been called
	newApp := NewOraichainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, uint(0), flags.FlagHome)
	_, _, err := newApp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func setGenesis(app *NewApp) error {
	genesisState := NewDefaultGenesisState()

	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			ChainId:       "Oraichain",
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	return nil
}
