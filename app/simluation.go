package app

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"

	aiRequest "github.com/oraichain/orai/x/airequest"
	aiResult "github.com/oraichain/orai/x/airesult"
	"github.com/oraichain/orai/x/provider"
	providerExported "github.com/oraichain/orai/x/provider/exported"
	providerTypes "github.com/oraichain/orai/x/provider/types"
	webSocket "github.com/oraichain/orai/x/websocket"
	"github.com/spf13/viper"

	tmdb "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// TestApp is used to test the whole Oraichain application. It is a wrapper that wraps around the NewApp struct.
type TestApp struct {
	NewApp
}

// Account is a mock account for testing
type Account struct {
	PrivKey    crypto.PrivKey
	PubKey     crypto.PubKey
	Address    sdk.AccAddress
	ValAddress sdk.ValAddress
}

// nolint
var (
	ChainID       = "Oraichain"
	DefaultCoins  = sdk.NewCoins(sdk.NewInt64Coin("orai", 100000000))
	TestCoins     = sdk.NewCoins(sdk.NewInt64Coin("orai", 99999999))
	MinimumFees   = sdk.NewCoins(sdk.NewInt64Coin("orai", 1000))
	MediumFees    = sdk.NewCoins(sdk.NewInt64Coin("orai", 10000))
	randSource    = int64(1111111111)
	keySize       = 100
	Duc           Account
	Bean          Account
	Bob           Account
	Ana           Account
	FirstVal      Account
	SecondVal     Account
	ThirdVal      Account
	OracleScripts []providerExported.OracleScriptI
	DataSources   []providerExported.AIDataSourceI
	TestCases     []providerExported.TestCaseI
)

// NewTestApp is the constructor of the TestApp struct
func NewTestApp() TestApp {
	// setup some prefix configurations for the testing app
	config := sdk.GetConfig()
	SetBech32AddressPrefixes(config)

	// Create a temp dir to run the simulation
	dir, err := ioutil.TempDir("", ".oraisimd")
	if err != nil {
		panic(err)
	}
	viper.Set(cli.HomeFlag, dir)

	// create a new instance of the OraichainApp for testing
	db := tmdb.NewMemDB()
	app := NewOraichainApp(log.NewNopLogger(), db, nil, true, uint(0), flags.FlagHome)

	// wrap the app with TestApp
	testApp := TestApp{NewApp: *app}

	// Initialize some default accounts for testing
	accs, valAccs := GenerateGenesisAccs()
	// Initialize the genesis state with some genesis validators based on the accounts gen above
	return testApp.InitializeFromGenesisStates(testApp.NewAuthGenState(accs, DefaultCoins), testApp.NewGenUtilState(valAccs, DefaultCoins))
}

// GenerateTestApp is used to create a simulation Oraichain application for unit tests
func GenerateTestApp() (TestApp, sdk.Context) {
	testApp := NewTestApp()
	//create a new context to have access to the KVStore and others
	ctx := testApp.NewContext(true, abci.Header{})

	return testApp, ctx
}

// InitScripts helps init some basic data sources, test cases and oracle scripts for testing
func (testApp TestApp) InitScripts(k provider.Keeper, ctx sdk.Context) {

	k.SetAIDataSource(ctx, "datasource", providerTypes.NewAIDataSource("datasource", Duc.Address, MinimumFees, "ABCDEF"))
	k.SetTestCase(ctx, "testcase", providerTypes.NewTestCase("testcase", Duc.Address, MinimumFees, "ABCDGH"))

	// // init scripts

	k.SetOracleScript(ctx, "testcase", providerTypes.NewOracleScript("oraclescript", Duc.Address, "AX", MediumFees, []string{"datasource"}, []string{"testcase"}))
}

// GetAccountKeeper is getter for the account keeper of TestApp
func (tApp TestApp) GetAccountKeeper() auth.AccountKeeper {
	return tApp.accountKeeper
}

// GetBankKeeper is getter for the bank keeper of TestApp
func (tApp TestApp) GetBankKeeper() bank.Keeper {
	return tApp.bankKeeper
}

// GetSupplyKeeper is getter for the supply keeper of TestApp
func (tApp TestApp) GetSupplyKeeper() supply.Keeper {
	return tApp.supplyKeeper
}

// GetStakingKeeper is getter for the staking keeper of TestApp
func (tApp TestApp) GetStakingKeeper() staking.Keeper {
	return tApp.stakingKeeper
}

// GetSlashingKeeper is getter for the slashing keeper of TestApp
func (tApp TestApp) GetSlashingKeeper() slashing.Keeper {
	return tApp.slashingKeeper
}

// GetMintKeeper is getter for the mint keeper of TestApp
func (tApp TestApp) GetMintKeeper() mint.Keeper {
	return tApp.mintKeeper
}

// GetDistrKeeper is getter for the distr keeper of TestApp
func (tApp TestApp) GetDistrKeeper() distribution.Keeper {
	return tApp.distrKeeper
}

// GetParamsKeeper is getter for the params keeper of TestApp
func (tApp TestApp) GetParamsKeeper() params.Keeper {
	return tApp.paramsKeeper
}

// GetProviderKeeper is getter for the provider keeper of TestApp
func (tApp TestApp) GetProviderKeeper() provider.Keeper {
	return tApp.providerKeeper
}

// GetAIRequestKeeper is getter for the airequest keeper of TestApp
func (tApp TestApp) GetAIRequestKeeper() aiRequest.Keeper {
	return tApp.aiRequestKeeper
}

// GetWebSocketKeeper is getter for the websocket keeper of TestApp
func (tApp TestApp) GetWebSocketKeeper() webSocket.Keeper {
	return tApp.webSocketKeeper
}

// GetAIResultKeeper is getter for the airesult keeper of TestApp
func (tApp TestApp) GetAIResultKeeper() aiResult.Keeper {
	return tApp.aiResultKeeper
}

// InitializeFromGenesisStates calls InitChain on the app using the default genesis state, overwitten with any passed in genesis states
func (tApp TestApp) InitializeFromGenesisStates(genesisStates ...GenesisState) TestApp {
	// Create a default genesis state and overwrite with provided values
	genesisState := NewDefaultGenesisState()
	for _, state := range genesisStates {
		for k, v := range state {
			genesisState[k] = v
		}
	}

	// create default genesis state for our modules
	providerGenesis := provider.DefaultGenesisState()
	aiRequestGenesis := aiRequest.DefaultGenesisState()
	webSocketGenesis := webSocket.DefaultGenesisState()
	aiResultGenesis := aiResult.DefaultGenesisState()

	// Add the default module genesis states into the app genesis state
	genesisState[provider.ModuleName] = tApp.Codec().MustMarshalJSON(providerGenesis)
	genesisState[aiRequest.ModuleName] = tApp.Codec().MustMarshalJSON(aiRequestGenesis)
	genesisState[webSocket.ModuleName] = tApp.Codec().MustMarshalJSON(webSocketGenesis)
	genesisState[aiResult.ModuleName] = tApp.Codec().MustMarshalJSON(aiResultGenesis)
	// Initialize the chain
	stateBytes, err := codec.MarshalJSONIndent(tApp.cdc, genesisState)
	if err != nil {
		panic(err)
	}
	// initiate the chain with minimum information, including the genesis state
	tApp.InitChain(
		abci.RequestInitChain{
			ChainId:       ChainID,
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	return tApp
}

// func (tApp TestApp) CheckBalance(t *testing.T, ctx sdk.Context, owner sdk.AccAddress, expectedCoins sdk.Coins) {
// 	acc := tApp.GetAccountKeeper().GetAccount(ctx, owner)
// 	require.NotNilf(t, acc, "account with address '%s' doesn't exist", owner)
// 	require.Equal(t, expectedCoins, acc.GetCoins())
// }

// NewAuthGenState creates a new auth genesis state for some test addresses and coins.
func (tApp TestApp) NewAuthGenState(accs []Account, coins sdk.Coins) GenesisState {
	// Create GenAccounts
	accounts := authexported.GenesisAccounts{}
	for i := range accs {
		accounts = append(accounts, auth.NewBaseAccount(accs[i].Address, coins, accs[i].PubKey, 0, 0))
	}
	// Create the auth genesis state
	authGenesis := auth.NewGenesisState(auth.DefaultParams(), accounts)
	return GenesisState{auth.ModuleName: tApp.Codec().MustMarshalJSON(authGenesis)}
}

// NewGenUtilState creates a new genutil genesis state to generate some genesis validators.
func (tApp TestApp) NewGenUtilState(vals []Account, coins sdk.Coins) GenesisState {
	// add auth types
	authTypes := []authtypes.StdTx{}
	for i := range vals {
		authTypes = append(authTypes, createValidatorTx(vals[i], fmt.Sprintf("%s %s ", "validator", strconv.Itoa(i)), DefaultCoins[0]))
	}
	// Create the genutil genesis state
	genutilGenesis := genutil.NewGenesisStateFromStdTx(authTypes)
	return GenesisState{genutil.ModuleName: tApp.Codec().MustMarshalJSON(genutilGenesis)}
}

// constructor function for the UserAccount struct
func generateAccount(r *rand.Rand) Account {
	privkeySeed := make([]byte, 32)
	_, err := r.Read(privkeySeed)
	if err != nil {
		panic("Could not read the random secret key")
	}
	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	return Account{
		PrivKey:    privKey,
		PubKey:     privKey.PubKey(),
		Address:    sdk.AccAddress(privKey.PubKey().Address()),
		ValAddress: sdk.ValAddress(privKey.PubKey().Address()),
	}
}

// GenerateGenesisAccs generate genesis accounts and validators for testing
func GenerateGenesisAccs() ([]Account, []Account) {
	// init user accounts for testing
	r := rand.New(rand.NewSource(randSource))
	Duc = generateAccount(r)
	Bean = generateAccount(r)
	Bob = generateAccount(r)
	Ana = generateAccount(r)
	// Validators also need accounts to interact with the system
	FirstVal = generateAccount(r)
	SecondVal = generateAccount(r)
	ThirdVal = generateAccount(r)

	var list []Account
	var valList []Account

	// add the user list and validators and return them to use
	list = append(list, Duc)
	list = append(list, Bean)
	list = append(list, Bob)
	list = append(list, Ana)
	list = append(list, FirstVal)
	list = append(list, SecondVal)
	list = append(list, ThirdVal)
	valList = append(valList, FirstVal)
	valList = append(valList, SecondVal)
	valList = append(valList, ThirdVal)

	return list, valList
}

// createValidatorTx is the genutil transaction to create genesis validator on the system
func createValidatorTx(acc Account, moniker string, selfDelegation sdk.Coin) authtypes.StdTx {
	msg := staking.NewMsgCreateValidator(
		acc.ValAddress, acc.PubKey, selfDelegation,
		staking.NewDescription(moniker, "", "", "", ""),
		staking.NewCommissionRates(sdk.MustNewDecFromStr("0.125"), sdk.MustNewDecFromStr("0.3"), sdk.MustNewDecFromStr("0.01")),
		sdk.NewInt(1),
	)

	// create a transaction to sign
	txMsg := authtypes.StdSignMsg{
		ChainID:       ChainID,
		AccountNumber: 0,
		Sequence:      0,
		Fee:           auth.NewStdFee(200000, sdk.Coins{}),
		Msgs:          []sdk.Msg{msg},
		Memo:          "",
	}
	// sign the message
	sigBytes, err := acc.PrivKey.Sign(txMsg.Bytes())
	if err != nil {
		panic(err)
	}
	sigs := []authtypes.StdSignature{{
		PubKey:    acc.PubKey,
		Signature: sigBytes,
	}}
	return authtypes.NewStdTx([]sdk.Msg{msg}, auth.NewStdFee(200000, sdk.Coins{}), sigs, "")
}
