package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/app"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

func TestSetAIRequest(t *testing.T) {
	// simulate the application
	testApp, ctx := app.GenerateTestApp()

	k := testApp.GetAIRequestKeeper()
	// get a random AI request ID
	id := ksuid.New().String()

	// collect a list of validators
	var listVal []sdk.ValAddress
	listVal = append(listVal, app.FirstVal.ValAddress)
	listVal = append(listVal, app.SecondVal.ValAddress)

	testApp.InitScripts(testApp.GetProviderKeeper(), ctx)

	k.SetAIRequest(ctx, id, types.NewAIRequest(id, "oscript", app.Duc.Address, listVal, ctx.BlockHeight(), app.DataSources, app.TestCases, app.MediumFees, []byte{}, []byte{}))

	_, err := k.GetAIRequest(ctx, id)

	require.NoError(t, err, "should not error because there is a data source named data source")
}
