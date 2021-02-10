package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAIDataSource is the constructor of the data source struct
func NewAIDataSource(
	name string,
	contract string,
	owner sdk.AccAddress,
	fees sdk.Coins,
	des string,
) *AIDataSource {
	return &AIDataSource{
		Name:        name,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
		Description: des,
	}
}

// verify interface at compile time
var _ sdk.Msg = &MsgCreateAIDataSource{}
var _ sdk.Msg = &MsgEditAIDataSource{}

// NewMsgCreateAIDataSource is a constructor function for MsgCreateAIDataSource
func NewMsgCreateAIDataSource(name string, contract string, owner sdk.AccAddress, fees string, des string) *MsgCreateAIDataSource {
	return &MsgCreateAIDataSource{
		Name:        name,
		Description: des,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
	}
}

// NewMsgEditAIDataSource is a constructor function for MsgEditAIDataSource
func NewMsgEditAIDataSource(oldName, newName string, contract string, owner sdk.AccAddress, fees string, des string) *MsgEditAIDataSource {
	return &MsgEditAIDataSource{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
	}
}

// NewOracleScript is the constructor of the oScript struct
func NewOracleScript(
	name string,
	contract string,
	owner sdk.AccAddress,
	des string,
	minimumFees sdk.Coins,
	dSources []string,
	tCases []string,
) *OracleScript {
	return &OracleScript{
		Name:        name,
		Contract:    contract,
		Owner:       owner,
		Description: des,
		MinimumFees: minimumFees,
		DSources:    dSources,
		TCases:      tCases,
	}
}

var _ sdk.Msg = &MsgCreateOracleScript{}
var _ sdk.Msg = &MsgEditOracleScript{}

// NewMsgCreateOracleScript is a constructor function for MsgCreateOracleScript
func NewMsgCreateOracleScript(name string, contract string, owner sdk.AccAddress, fees string, des string, dSources, tCases []string) *MsgCreateOracleScript {
	return &MsgCreateOracleScript{
		Name:        name,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
		Description: des,
		DataSources: dSources,
		TestCases:   tCases,
	}
}

// NewMsgEditOracleScript is a constructor function for MsgEditOracleScript
func NewMsgEditOracleScript(oldName, newName string, contract string, owner sdk.AccAddress, fees string, des string, dSources, tCases []string) *MsgEditOracleScript {
	return &MsgEditOracleScript{
		OldName:     oldName,
		NewName:     newName,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
		Description: des,
		DataSources: dSources,
		TestCases:   tCases,
	}
}

// NewTestCase is the constructor of the testcase struct
func NewTestCase(
	name string,
	contract string,
	owner sdk.AccAddress,
	fees sdk.Coins,
	des string,
) *TestCase {
	return &TestCase{
		Name:        name,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
		Description: des,
	}
}

// verify interface at compile time
var _ sdk.Msg = &MsgCreateTestCase{}
var _ sdk.Msg = &MsgEditTestCase{}

// NewMsgCreateTestCase is a constructor function for MsgCreateTestCase
func NewMsgCreateTestCase(name string, contract string, owner sdk.AccAddress, fees string, des string) *MsgCreateTestCase {
	return &MsgCreateTestCase{
		Name:        name,
		Description: des,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
	}
}

// NewMsgEditTestCase is a constructor function for MsgEditTestCase
func NewMsgEditTestCase(oldName, newName string, contract string, owner sdk.AccAddress, fees string, des string) *MsgEditTestCase {
	return &MsgEditTestCase{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Contract:    contract,
		Owner:       owner,
		Fees:        fees,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(oScripts []OracleScript, aiDataSources []AIDataSource, testCases []TestCase) *GenesisState {
	return &GenesisState{
		// TODO: Fill out according to your genesis state
		OracleScripts: oScripts,
		AIDataSources: aiDataSources,
		TestCases:     testCases,
	}
}
