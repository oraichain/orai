package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/websocket/exported"
)

// MsgCreateReport defines message for creating a report by a reporter of a validator
type MsgCreateReport struct {
	RequestID         string                       `json:"request_id"`
	DataSourceResults []exported.DataSourceResultI `json:"data_source_results"`
	TestCaseResults   []exported.TestCaseResultI   `json:"test_case_results"`
	Reporter          Reporter                     `json:"reporter"`
	Fees              sdk.Coins                    `json:"report_fee"`
	AggregatedResult  []byte                       `json:"aggregated_result"`
	ResultStatus      string                       `json:"result_status"`
}

// NewMsgCreateReport is a constructor function for MsgCreateReport
func NewMsgCreateReport(
	requestID string,
	dataSourceResults []exported.DataSourceResultI,
	testCaseResults []exported.TestCaseResultI,
	reporter Reporter,
	fees sdk.Coins,
	aggregatedResult []byte,
	status string,
) MsgCreateReport {
	return MsgCreateReport{
		RequestID:         requestID,
		DataSourceResults: dataSourceResults,
		TestCaseResults:   testCaseResults,
		Reporter:          reporter,
		Fees:              fees,
		AggregatedResult:  aggregatedResult,
		ResultStatus:      status,
	}
}

// Route should return the name of the module
func (msg MsgCreateReport) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateReport) Type() string { return "create_report" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateReport) ValidateBasic() error {
	if msg.Reporter.Address.Empty() || len(msg.Reporter.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Reporter.String())
	} else if len(msg.RequestID) == 0 || msg.Reporter.Validator.Empty() {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "Request ID / validator address cannot be empty")
	} else if len(msg.DataSourceResults) == 0 || len(msg.TestCaseResults) == 0 || len(msg.AggregatedResult) == 0 {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "lengths of the data source and test case must be greater than zero, and there must be an aggregated result")
	} else if msg.ResultStatus != ResultSuccess && msg.ResultStatus != ResultFailure {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "result status of the report is not valid")
	} else {
		_, err := sdk.ParseCoins(msg.Fees.String())
		if err != nil {
			return sdkerrors.Wrap(ErrReportFeeTypeInvalid, err.Error())
		}
		return nil
	}
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateReport) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateReport) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Reporter.Address}
}
