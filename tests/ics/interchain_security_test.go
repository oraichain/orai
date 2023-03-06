package ics

import (
	"encoding/json"
	"fmt"
	"testing"

	appConsumer "github.com/cosmos/interchain-security/app/consumer"
	"github.com/cosmos/interchain-security/tests/e2e"
	icstestingutils "github.com/cosmos/interchain-security/testutil/ibc_testing"
	"github.com/stretchr/testify/suite"

	"github.com/CosmWasm/wasmd/x/wasm"
	ibctesting "github.com/cosmos/interchain-security/legacy_ibc_testing/testing"
	oraiApp "github.com/oraichain/orai/app"
	"github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"
)

func TestCCVTestSuite(t *testing.T) {
	// Pass in concrete app types that implement the interfaces defined in https://github.com/cosmos/interchain-security/testutil/e2e/interfaces.go
	// IMPORTANT: the concrete app types passed in as type parameters here must match the
	// concrete app types returned by the relevant app initers.
	fmt.Println("run test case: TestCCVTestSuite")
	ccvSuite := e2e.NewCCVTestSuite[*oraiApp.OraichainApp, *appConsumer.App](
		// Pass in ibctesting.AppIniters for orai (provider) and consumer.
		OraiAppIniter, icstestingutils.ConsumerAppIniter, []string{})

	// Run tests
	suite.Run(t, ccvSuite)
}

// OraiAppIniter implements ibctesting.AppIniter for the orai app
func OraiAppIniter() (ibctesting.TestingApp, map[string]json.RawMessage) {
	encoding := oraiApp.MakeEncodingConfig()
	app := oraiApp.NewOraichainApp(log.NewNopLogger(), tmdb.NewMemDB(), nil, true, map[int64]bool{},
		oraiApp.DefaultNodeHome, 5, encoding, wasm.EnableAllProposals, oraiApp.EmptyAppOptions{}, nil)
	testApp := ibctesting.TestingApp(app)
	return testApp, oraiApp.NewDefaultGenesisState(app.AppCodec())
}
