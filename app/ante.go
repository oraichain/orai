package app

import (
	"fmt"
	"runtime/debug"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	ibcante "github.com/cosmos/ibc-go/v4/modules/core/ante"
	"github.com/cosmos/ibc-go/v4/modules/core/keeper"
	customante "github.com/oraichain/orai/app/ante"
	tmlog "github.com/tendermint/tendermint/libs/log"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evmante "github.com/tharsis/ethermint/app/ante"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	globalfeeante "github.com/CosmosContracts/juno/v15/x/globalfee/ante"
	globalfeekeeper "github.com/CosmosContracts/juno/v15/x/globalfee/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const maxBypassMinFeeMsgGasUsage = 1_000_000

// HandlerOptions extends the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	AccountKeeper  evmtypes.AccountKeeper
	BankKeeper     evmtypes.BankKeeper
	EvmKeeper      evmante.EVMKeeper
	FeegrantKeeper ante.FeegrantKeeper

	GlobalFeeKeeper globalfeekeeper.Keeper
	StakingKeeper   stakingkeeper.Keeper

	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  ante.SignatureVerificationGasConsumer
	FeeMarketKeeper evmtypes.FeeMarketKeeper
	MaxTxGasWanted  uint64

	IBCKeeper         *keeper.Keeper
	TxCounterStoreKey sdk.StoreKey
	WasmConfig        wasmTypes.WasmConfig
	Cdc               codec.BinaryCodec

	BypassMinFeeMsgTypes []string
}

func (options HandlerOptions) Validate() error {
	if options.AccountKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.EvmKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "evm keeper is required for AnteHandler")
	}
	if options.IBCKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "ibc keeper is required for AnteHandler")
	}
	if options.FeeMarketKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee market keeper is required for AnteHandler")
	}
	if options.EvmKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "evm keeper is required for AnteHandler")
	}
	if options.SigGasConsumer == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "signature gas consumer is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for AnteHandler")
	}
	return nil
}

// cosmosHandlerOptions extends HandlerOptions to provide some Cosmos specific configurations
type cosmosHandlerOptions struct {
	HandlerOptions
	isEIP712 bool
}

// NewAnteHandler returns an 'AnteHandler' that will run actions before a tx is sent to a module's handler.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		defer Recover(ctx.Logger(), &err)

		txWithExtensions, ok := tx.(ante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 1 {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest,
					"rejecting tx with more than 1 extension option",
				)
			}

			if len(opts) == 1 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
					// handle as *evmtypes.MsgEthereumTx
					anteHandler = newEthAnteHandler(options)
				case "/ethermint.types.v1.ExtensionOptionsWeb3Tx":
					// handle as normal Cosmos SDK tx, except signature is checked for EIP712 representation
					anteHandler = newCosmosAnteHandler(cosmosHandlerOptions{
						HandlerOptions: options,
						isEIP712:       true,
					})
				default:
					return ctx, sdkerrors.Wrapf(
						sdkerrors.ErrUnknownExtensionOptions,
						"rejecting tx with unsupported extension option: %s", typeURL,
					)
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		// handle as totally normal Cosmos SDK tx
		switch tx.(type) {
		case sdk.Tx:
			anteHandler = newCosmosAnteHandler(cosmosHandlerOptions{
				HandlerOptions: options,
				isEIP712:       false,
				
			})
		default:
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}, nil
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func newCosmosAnteHandler(options cosmosHandlerOptions) sdk.AnteHandler {

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	decorators := []sdk.AnteDecorator{}

	decorators = append(decorators,
		evmante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		ante.NewSetUpContextDecorator(),   // second decorator. SetUpContext must be called before other decorators
	)

	if !options.isEIP712 {
		decorators = append(decorators, ante.NewRejectExtensionOptionsDecorator())
	}

	var sigVerification sdk.AnteDecorator = ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler)
	if options.isEIP712 {
		sigVerification = evmante.NewEip712SigVerificationDecorator(options.AccountKeeper, options.SignModeHandler, options.EvmKeeper)
	}

	decorators = append(decorators,
		customante.NewEvmMinGasFilter(options.EvmKeeper), // filter out evm denom from min-gas-prices
		ante.NewMempoolFeeDecorator(),
		customante.NewVestingAccountDecorator(),
		customante.NewMinCommissionDecorator(options.Cdc),
		customante.NewAuthzLimiterDecorator(
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
			sdk.MsgTypeURL(&vesting.MsgCreateVestingAccount{}),
		),
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit),
		wasmkeeper.NewCountTXDecorator(options.TxCounterStoreKey),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		globalfeeante.NewFeeDecorator(options.BypassMinFeeMsgTypes, options.GlobalFeeKeeper, options.StakingKeeper, maxBypassMinFeeMsgGasUsage),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		// SetPubKeyDecorator must be called before all signature verification decorators
		evmante.NewSetPubKeyDecorator(options.AccountKeeper, options.EvmKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		sigVerification,
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewAnteDecorator(options.IBCKeeper),
	)

	return sdk.ChainAnteDecorators(decorators...)
}

func newEthAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		evmante.NewEthSetUpContextDecorator(options.EvmKeeper), // outermost AnteDecorator. SetUpContext must be called first
		evmante.NewEthMempoolFeeDecorator(options.EvmKeeper),   // Check eth effective gas price against minimal-gas-prices
		evmante.NewEthValidateBasicDecorator(options.EvmKeeper),
		evmante.NewEthSigVerificationDecorator(options.EvmKeeper),
		evmante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.BankKeeper, options.EvmKeeper),
		evmante.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
		evmante.NewCanTransferDecorator(options.EvmKeeper),
		evmante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper, options.EvmKeeper), // innermost AnteDecorator.
	)
}

func Recover(logger tmlog.Logger, err *error) {
	if r := recover(); r != nil {
		*err = sdkerrors.Wrapf(sdkerrors.ErrPanic, "%v", r)

		if e, ok := r.(error); ok {
			logger.Error(
				"ante handler panicked",
				"error", e,
				"stack trace", string(debug.Stack()),
			)
		} else {
			logger.Error(
				"ante handler panicked",
				"recover", fmt.Sprintf("%v", r),
			)
		}
	}
}
