package keeper_test

import (
	"testing"

	"github.com/oraichain/orai/app"
	"github.com/oraichain/orai/x/provider/types"
	"github.com/stretchr/testify/require"
)

func TestHasOracleScript(t *testing.T) {
	testApp, ctx := app.GenerateTestApp()
	testApp.GetProviderKeeper().SetOracleScript(ctx, "abc", types.NewOracleScript("abc", app.Duc.Address, "Hello there this is a test oracle script", app.MinimumFees, []string{"aa"}, []string{"ccc"}))
	_, err := testApp.GetProviderKeeper().GetOracleScript(ctx, "abc")
	require.Error(t, err, "should receive error because there is an oracle script named abc")
}
