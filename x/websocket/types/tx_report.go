package types

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// regex allow only alphabet, numeric and underscore characters
var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString

// Route should return the name of the module
func (msg *MsgCreateReport) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCreateReport) Type() string { return "create_report" }

// ValidateBasic runs stateless checks on the message
func (msg *MsgCreateReport) ValidateBasic() error {
	reporter := msg.GetReporter()
	if reporter.GetAddress().Empty() || len(reporter.GetName()) == 0 || isStringAlphabetic(reporter.GetName()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, reporter.String())
	} else if len(msg.GetRequestID()) == 0 || reporter.Validator.Empty() {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "Request ID / validator address cannot be empty")
	} else if len(msg.GetDataSourceResults()) == 0 || len(msg.GetTestCaseResults()) == 0 || len(msg.GetAggregatedResult()) == 0 {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "lengths of the data source and test case must be greater than zero, and there must be an aggregated result")
	} else if msg.GetResultStatus() != ResultSuccess && msg.GetResultStatus() != ResultFailure {
		return sdkerrors.Wrap(ErrMsgReportInvalid, "result status of the report is not valid")
	} else {
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
