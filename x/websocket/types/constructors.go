package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// verify interface at compile time
var _ sdk.Msg = &MsgCreateReport{}
var _ sdk.Msg = &MsgAddReporter{}
var _ sdk.Msg = &MsgRemoveReporter{}

// NewMsgCreateReport is a constructor function for MsgCreateReport
func NewMsgCreateReport(
	requestID string,
	dataSourceResults []*DataSourceResult,
	testCaseResults []*TestCaseResult,
	reporter *Reporter,
	fees sdk.Coins,
	aggregatedResult []byte,
	status string,
) *MsgCreateReport {
	return &MsgCreateReport{
		RequestID:         requestID,
		DataSourceResults: dataSourceResults,
		TestCaseResults:   testCaseResults,
		Reporter:          reporter,
		Fees:              fees,
		AggregatedResult:  aggregatedResult,
		ResultStatus:      status,
	}
}

// NewMsgAddReporter is a constructor function for MsgAddReporter
func NewMsgAddReporter(validator sdk.ValAddress, reporter sdk.AccAddress, adder sdk.AccAddress) *MsgAddReporter {
	return &MsgAddReporter{
		Adder:     adder,
		Validator: validator,
		Reporter:  reporter,
	}
}

// NewMsgRemoveReporter is a constructor function for MsgRemoveReporter
func NewMsgRemoveReporter(validator sdk.ValAddress, reporter sdk.AccAddress, remover sdk.AccAddress) *MsgRemoveReporter {
	return &MsgRemoveReporter{
		Remover:   remover,
		Validator: validator,
		Reporter:  reporter,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(reports []Report, reporters []Reporter) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		Reports:   reports,
		Reporters: reporters,
	}
}
