package keeper_test

import (
	"testing"

	"github.com/oraichain/orai/app"
	"github.com/oraichain/orai/x/provider/types"
	"github.com/stretchr/testify/require"
)

func TestSetDataSource(t *testing.T) {
	// simulate the application
	testApp, ctx := app.GenerateTestApp()

	k := testApp.GetProviderKeeper()
	k.SetAIDataSource(ctx, "datasource2", types.NewAIDataSource("datasource2", app.Duc.Address, app.MinimumFees, "ABCD"))

	_, err := k.GetAIDataSource(ctx, "datasource2")

	require.NoError(t, err, "should not error because there is a data source named data source")
}
