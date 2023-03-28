package e2e_test

import (
	"encoding/json"
	"testing"

	appConsumer "github.com/cosmos/interchain-security/app/consumer"
	"github.com/cosmos/interchain-security/tests/e2e"
	icstestingutils "github.com/cosmos/interchain-security/testutil/ibc_testing"
	"github.com/stretchr/testify/suite"

	ibctesting "github.com/cosmos/interchain-security/legacy_ibc_testing/testing"

	oraiApp "github.com/oraichain/orai/app"
	"github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"
)

// OraiAppIniter implements ibctesting.AppIniter for the orai app
func OraiAppIniter() (ibctesting.TestingApp, map[string]json.RawMessage) {
	encoding := oraiApp.MakeEncodingConfig()
	app := oraiApp.NewOraichainApp(log.NewNopLogger(), tmdb.NewMemDB(), nil, true, map[int64]bool{},
		oraiApp.DefaultNodeHome, 5, encoding, oraiApp.GetEnabledProposals(), oraiApp.EmptyAppOptions{}, nil)
	testApp := ibctesting.TestingApp(app)
	return testApp, oraiApp.ModuleBasics.DefaultGenesis(encoding.Marshaler)
}

// Executes the standard group of ccv tests against a consumer and provider app.go implementation.
func TestCCVTestSuite(t *testing.T) {
	// Pass in concrete app types that implement the interfaces defined in /testutil/e2e/interfaces.go
	ccvSuite := e2e.NewCCVTestSuite[*oraiApp.OraichainApp, *appConsumer.App](
		// Pass in ibctesting.AppIniters for provider and consumer.
		OraiAppIniter, icstestingutils.ConsumerAppIniter,
		// TODO: These three tests just don't work in IS, so skip them for now
		[]string{})

	// Run tests
	suite.Run(t, ccvSuite)
}
