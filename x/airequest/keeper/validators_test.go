package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/packages/rng"
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
)

func TestCalucateMol(t *testing.T) {
	k := &keeper.Keeper{}
	size := 5
	maxValidatorSize := 5
	totalPowers := int64(1000000)
	randomGenerator, _ := rng.NewRng(make([]byte, types.RngSeedSize), []byte("nonce"), []byte("Oraichain"))
	valOperators := make([]sdk.ValAddress, maxValidatorSize)
	t.Logf("Validator length: %v\n", len(valOperators))
	for i := 0; i < maxValidatorSize; i++ {
		valOperators[i] = ed25519.GenPrivKey().PubKey().Address().Bytes()
	}

	validators := k.SampleIndexes(valOperators, size, randomGenerator, totalPowers)
	t.Logf("Max Validators :%v\n", valOperators)
	t.Logf("Choosen Validators :%v\n", validators)
}
