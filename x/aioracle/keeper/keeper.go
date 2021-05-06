package keeper

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/oraichain/orai/x/aioracle/types"
	"github.com/tendermint/tendermint/libs/log"
)

// always clone keeper to make it immutable
type (
	Keeper struct {
		Cdc              codec.Marshaler
		StoreKey         sdk.StoreKey
		WasmKeeper       *wasm.Keeper
		ParamSpace       params.Subspace
		StakingKeeper    staking.Keeper
		BankKeeper       bank.Keeper
		DistrKeeper      distr.Keeper
		AuthKeeper       auth.AccountKeeper
		FeeCollectorName string
	}
)

// NewKeeper creates a aioracle keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, wasmKeeper *wasm.Keeper, aiOracleSubspace params.Subspace, stakingKeeper staking.Keeper, bankKeeper bank.Keeper, authKeeper auth.AccountKeeper, distrKeeper distr.Keeper, feeCollector string) *Keeper {
	if !aiOracleSubspace.HasKeyTable() {
		// register parameters of the aioracle module into the param space
		aiOracleSubspace = aiOracleSubspace.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		StoreKey:         key,
		Cdc:              cdc,
		WasmKeeper:       wasmKeeper,
		ParamSpace:       aiOracleSubspace,
		StakingKeeper:    stakingKeeper,
		BankKeeper:       bankKeeper,
		DistrKeeper:      distrKeeper,
		AuthKeeper:       authKeeper,
		FeeCollectorName: feeCollector,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// NewValResult is a wrapper function of the aioracle module that allow others to initiate a new valresult entity through the keeper
func (k *Keeper) NewValResult(val *types.Validator, result []byte, status string) *types.ValResult {
	return &types.ValResult{
		Validator:    val,
		Result:       result,
		ResultStatus: status,
	}
}

// GetValidator return a specific validator given a validator address
func (k *Keeper) GetValidator(ctx sdk.Context, valAddress sdk.ValAddress) stakingtypes.ValidatorI {
	return k.StakingKeeper.Validator(ctx, valAddress)
}

// NewValidator is a wrapper function of the aioracle module that allow others to initiate a new validator entity through the keeper
func (k *Keeper) NewValidator(address sdk.ValAddress, votingPower int64, status string) *types.Validator {
	return &types.Validator{
		Address:     address,
		VotingPower: votingPower,
		Status:      status,
	}
}

// QueryContract return data from input of smart contract, should be struct with json serialized
func (k *Keeper) QueryContract(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	return k.WasmKeeper.QuerySmart(ctx, contractAddr, req)
}

func (k *Keeper) CalculateValidatorFees(ctx sdk.Context, providerFees sdk.Coins) sdk.Coins {
	// change reward ratio to the ratio of validator
	rewardRatio := k.GetParam(ctx, types.KeyAIOracleRewardPercentages)
	//rewardRatio := 40
	// reward = 1 - oracle reward percentage × test case fees
	valFees, _ := sdk.NewDecCoinsFromCoins(providerFees...).MulDec(sdk.NewDecWithPrec(int64(rewardRatio), 2)).TruncateDecimal()
	return valFees
}
