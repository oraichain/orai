package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAIDataSource is the constructor of the data source struct
func NewAIDataSource(
	name string,
	owner sdk.AccAddress,
	fees sdk.Coins,
	des string,
) AIDataSource {
	return AIDataSource{
		Name:        name,
		Owner:       owner,
		Fees:        fees,
		Description: des,
	}
}

// NewMsgCreateAIDataSource is a constructor function for MsgCreateAIDataSource
func NewMsgCreateAIDataSource(name string, code []byte, owner sdk.AccAddress, fees string, des string) MsgCreateAIDataSource {
	return MsgCreateAIDataSource{
		Name:        name,
		Description: des,
		Code:        code,
		Owner:       owner,
		Fees:        fees,
	}
}

// NewOracleScript is the constructor of the oScript struct
func NewOracleScript(
	name string,
	owner sdk.AccAddress,
	des string,
	minimumFees sdk.Coins,
	dSources []string,
	tCases []string,
) OracleScript {
	return OracleScript{
		Name:        name,
		Owner:       owner,
		Description: des,
		MinimumFees: minimumFees,
		DSources:    dSources,
		TCases:      tCases,
	}
}

// NewMsgCreateOracleScript is a constructor function for MsgCreateOracleScript
func NewMsgCreateOracleScript(name string, code []byte, owner sdk.AccAddress, des string, dSources, tCases []string) MsgCreateOracleScript {
	return MsgCreateOracleScript{
		Name:        name,
		Code:        code,
		Owner:       owner,
		Description: des,
		DataSources: dSources,
		TestCases:   tCases,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(oScripts []OracleScript, aiDataSources []AIDataSource, testCases []TestCase) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		OracleScripts: oScripts,
		AIDataSources: aiDataSources,
		TestCases:     testCases,
	}
}
