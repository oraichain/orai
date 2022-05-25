package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	provider "github.com/oraichain/orai/x/provider/types"
)

// Route should return the name of the module
func (msg *MsgCreateReport) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCreateReport) Type() string { return "create_report" }

// ValidateBasic runs stateless checks on the message
func (msg *MsgCreateReport) ValidateBasic() error {
	reporter := msg.GetReporter()
	if reporter.GetAddress().Empty() || len(reporter.GetName()) == 0 || !provider.IsStringAlphabetic(reporter.GetName()) || len(reporter.GetName()) >= ReporterNameLen {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, reporter.String())
	} else if len(msg.GetRequestID()) == 0 || reporter.Validator.Empty() {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "Request ID / validator address cannot be empty")
	} else if len(msg.GetDataSourceResults()) == 0 || len(msg.GetTestCaseResults()) == 0 || len(msg.GetAggregatedResult()) == 0 {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "lengths of the data source and test case must be greater than zero, and there must be an aggregated result")
	} else if msg.GetResultStatus() != ResultSuccess && msg.GetResultStatus() != ResultFailure {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "result status of the report is not valid")
	} else {
		var dsResultSize int
		for _, dsResult := range msg.DataSourceResults {
			dsResultSize += len(dsResult.Result)
		}
		var tcResultSize int
		for _, tcResult := range msg.TestCaseResults {
			for _, dsResult := range tcResult.DataSourceResults {
				tcResultSize += len(dsResult.Result)
			}
		}
		aggregatedResultSize := len(msg.AggregatedResult)
		requestIdSize := len(msg.RequestID)
		finalLen := dsResultSize + tcResultSize + aggregatedResultSize + requestIdSize
		if finalLen >= MsgLen {
			return sdkerrors.Wrap(ErrMsgReportInvalid, "Size of the report should not be larger than 200KB")
		}

		_, err := sdk.ParseCoinsNormalized(msg.Fees.String())
		if err != nil {
			return sdkerrors.Wrap(ErrReportFeeTypeInvalid, err.Error())
		}
		return nil
	}
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateReport) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgCreateReport) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Reporter.Address}
}
