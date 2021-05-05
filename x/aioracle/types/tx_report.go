package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgCreateReport) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCreateReport) Type() string { return EventTypeSetReport }

// ValidateBasic runs stateless checks on the message
func (msg *MsgCreateReport) ValidateBasic() error {
	if len(msg.GetDataSourceResults()) <= 0 || len(msg.GetAggregatedResult()) <= 0 {
		return fmt.Errorf("Report results are invalid")
	}
	if msg.GetResultStatus() != ResultFailure && msg.GetResultStatus() != ResultSuccess {
		return fmt.Errorf("Report result status is invalid: %v", msg.GetResultStatus())
	}
	for _, dsResult := range msg.GetDataSourceResults() {
		if dsResult.GetStatus() != ResultFailure && dsResult.GetStatus() != ResultSuccess {
			return fmt.Errorf("Data source result status is invalid: %v", dsResult.GetResult())
		}
	}
	if len(msg.BaseReport.GetValidatorAddress()) <= 0 || len(msg.BaseReport.GetRequestId()) <= 0 || msg.BaseReport.GetBlockHeight() <= 0 {
		return fmt.Errorf("Data source result status is invalid: %v", msg)
	}
	_, err := sdk.ParseCoinsNormalized(msg.BaseReport.Fees.String())
	if err != nil {
		return sdkerrors.Wrap(ErrReportFeeTypeInvalid, err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateReport) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgCreateReport) GetSigners() []sdk.AccAddress {
	creator := sdk.AccAddress(msg.BaseReport.ValidatorAddress)
	return []sdk.AccAddress{creator}
}

// Route should return the name of the module
func (msg *MsgCreateTestCaseReport) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgCreateTestCaseReport) Type() string { return EventTypeSetTestCaseReport }

// ValidateBasic runs stateless checks on the message
func (msg *MsgCreateTestCaseReport) ValidateBasic() error {
	if len(msg.GetResultsWithTestCase()) <= 0 {
		return fmt.Errorf("Report results are invalid")
	}
	for _, result := range msg.GetResultsWithTestCase() {
		if result.GetStatus() != ResultFailure && result.GetStatus() != ResultSuccess {
			return fmt.Errorf("Result status is invalid: %v", result)
		}
		for _, tcResult := range result.GetTestCaseResults() {
			if tcResult.GetStatus() != ResultFailure && tcResult.GetStatus() != ResultSuccess {
				return fmt.Errorf("Report result status is invalid: %v", tcResult.GetStatus())
			}
		}
	}
	if len(msg.BaseReport.GetValidatorAddress()) <= 0 || len(msg.BaseReport.GetRequestId()) <= 0 || msg.BaseReport.GetBlockHeight() <= 0 {
		return fmt.Errorf("Data source result status is invalid: %v", msg)
	}
	_, err := sdk.ParseCoinsNormalized(msg.BaseReport.Fees.String())
	if err != nil {
		return sdkerrors.Wrap(ErrReportFeeTypeInvalid, err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateTestCaseReport) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgCreateTestCaseReport) GetSigners() []sdk.AccAddress {
	creator := sdk.AccAddress(msg.BaseReport.ValidatorAddress)
	return []sdk.AccAddress{creator}
}
