package keeper_test

import (
	gocontext "context"
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/oraichain/orai/x/provider/types"
)

func TestQueryCurrentPlan(t *testing.T) {
	network := network.New(t, network.DefaultConfig())
	_, err := network.WaitForHeight(2)
	t.Log(err)
	val0 := network.Validators[0]
	queryClient := types.NewQueryClient(val0.ClientCtx)

	req := &types.DataSourceInfoReq{}

	res, err := queryClient.DataSourceInfo(gocontext.Background(), req)
	t.Log(res, err)

}
