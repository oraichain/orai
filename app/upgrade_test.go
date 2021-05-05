package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpgradeStoreLoader is used to prepare baseapp with a fixed StoreLoader
// pattern. This is useful for custom upgrade loading logic.
func upgradeStoreLoader(upgradeHeight int64, storeUpgrades *store.StoreUpgrades) baseapp.StoreLoader {
	return func(ms sdk.CommitMultiStore) error {
		if upgradeHeight == ms.LastCommitID().Version+1 {
			// Check if the current commit version and upgrade height matches
			if len(storeUpgrades.Renamed) > 0 || len(storeUpgrades.Deleted) > 0 || len(storeUpgrades.Added) > 0 {
				return ms.LoadLatestVersionAndUpgrade(storeUpgrades)
			}
		}

		// Otherwise load default store loader
		return baseapp.DefaultStoreLoader(ms)
	}
}

func useUpgradeLoader(height int64, upgrades *store.StoreUpgrades) func(*baseapp.BaseApp) {
	return func(app *baseapp.BaseApp) {
		app.SetStoreLoader(upgradeStoreLoader(height, upgrades))
	}
}

func defaultLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
}

func newStore(db dbm.DB, storeKeys []string) (*rootmulti.Store, map[string]*store.KVStoreKey, error) {
	rs := rootmulti.NewStore(db)
	rs.SetPruning(store.PruneNothing)

	keys := sdk.NewKVStoreKeys(storeKeys...)
	for _, key := range keys {
		rs.MountStoreWithDB(key, store.StoreTypeIAVL, nil)
	}
	err := rs.LoadLatestVersion()
	return rs, keys, err
}

func initStore(t *testing.T, db dbm.DB, storeKeys []string, k, v []byte) {
	rs, keys, err := newStore(db, storeKeys)
	require.Nil(t, err)
	require.Equal(t, int64(0), rs.LastCommitID().Version)

	// write some data in substore
	kv, _ := rs.GetStore(keys[storeKeys[0]]).(store.KVStore)
	require.NotNil(t, kv)
	kv.Set(k, v)
	commitID := rs.Commit()
	require.Equal(t, int64(1), commitID.Version)
}

func checkStore(t *testing.T, db dbm.DB, ver int64, storeKeys []string, k, v []byte) {
	rs, keys, err := newStore(db, storeKeys)
	require.Nil(t, err)
	require.Equal(t, ver, rs.LastCommitID().Version)

	// query data in substore
	kv, _ := rs.GetStore(keys[storeKeys[0]]).(store.KVStore)
	fmt.Println("store: ", keys, kv)
	fmt.Println("value: ", string(kv.Get(k)))
	require.Equal(t, string(v), string(kv.Get(k)))
}

// Test that we can make commits and then reload old versions.
// Test that LoadLatestVersion actually does.
func TestSetLoader(t *testing.T) {
	upgradeHeight := int64(5)

	// set a temporary home dir
	homeDir := t.TempDir()
	upgradeInfoFilePath := filepath.Join(homeDir, "upgrade-info.json")
	upgradeInfo := &store.UpgradeInfo{
		Name: "test", Height: upgradeHeight,
	}

	data, err := json.Marshal(upgradeInfo)
	require.NoError(t, err)

	err = ioutil.WriteFile(upgradeInfoFilePath, data, 0644)
	require.NoError(t, err)

	// make sure it exists before running everything
	_, err = os.Stat(upgradeInfoFilePath)
	require.NoError(t, err)

	cases := map[string]struct {
		setLoader    func(*baseapp.BaseApp)
		origStoreKey []string
		loadStoreKey []string
	}{
		"rename with inline opts": {
			setLoader: useUpgradeLoader(upgradeHeight, &store.StoreUpgrades{
				Deleted: []string{"provider", "aioracle", "websocket", "airesult"},
				// Renamed: []store.StoreRename{{
				// 	OldKey: "provider",
				// 	NewKey: "test",
				// }},
				Added: []string{"aioracle"},
			}),
			origStoreKey: []string{"provider", "aioracle", "websocket", "airesult"},
			loadStoreKey: []string{"aioracle"},
		},
	}

	k := []byte("key")
	v := []byte("value")

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			// prepare a db with some data
			db := dbm.NewMemDB()

			initStore(t, db, tc.origStoreKey, k, v)

			// load the app with the existing db
			opts := []func(*baseapp.BaseApp){baseapp.SetPruning(store.PruneNothing)}

			origapp := baseapp.NewBaseApp(t.Name(), defaultLogger(), db, nil, opts...)
			originalKeys := sdk.NewKVStoreKeys(tc.origStoreKey...)
			origapp.MountKVStores(originalKeys)
			err := origapp.LoadLatestVersion()
			require.Nil(t, err)

			for i := int64(2); i <= upgradeHeight-1; i++ {
				origapp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: i}})
				res := origapp.Commit()
				require.NotNil(t, res.Data)
			}

			if tc.setLoader != nil {
				opts = append(opts, tc.setLoader)
			}

			// load the new app with the original app db
			app := baseapp.NewBaseApp(t.Name(), defaultLogger(), db, nil, opts...)
			storeKeys := sdk.NewKVStoreKeys(tc.loadStoreKey...)
			app.MountKVStores(storeKeys)
			err = app.LoadLatestVersion()
			require.Nil(t, err)
			app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: upgradeHeight}})
			res := app.Commit()
			require.NotNil(t, res.Data)

			rs, keys, err := newStore(db, tc.loadStoreKey)
			// query data in substore
			kv, _ := rs.GetStore(keys["aioracle"]).(store.KVStore)
			kv.Set(k, v)
			require.NotNil(t, kv)
			rs.Commit()

			// check db is properly updated
			checkStore(t, db, upgradeHeight+1, tc.loadStoreKey, k, v)
		})
	}
}
