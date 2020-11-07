package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetSigners defines whose signature is required
func (msg MsgSetPriceRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}

// MsgCreateReport defines message for creating a report by a reporter of a validator
type MsgCreateReport struct {
	RequestID         string             `json:"request_id"`
	Validator         sdk.ValAddress     `json:"validator"`
	DataSourceResults []DataSourceResult `json:"data_source_results"`
	TestCaseResults   []TestCaseResult   `json:"test_case_results"`
	Reporter          sdk.AccAddress     `json:"reporter"`
	Fees              sdk.Coins          `json:"report_fee"`
	AggregatedResult  []byte             `json:"aggregated_result"`
}

// NewMsgCreateReport is a constructor function for MsgCreateReport
func NewMsgCreateReport(
	requestID string,
	validator sdk.ValAddress,
	dataSourceResults []DataSourceResult,
	testCaseResults []TestCaseResult,
	reporter sdk.AccAddress,
	fees sdk.Coins,
	aggregatedResult []byte,
) MsgCreateReport {
	return MsgCreateReport{
		RequestID:         requestID,
		Validator:         validator,
		DataSourceResults: dataSourceResults,
		TestCaseResults:   testCaseResults,
		Reporter:          reporter,
		Fees:              fees,
		AggregatedResult:  aggregatedResult,
	}
}

// Route should return the name of the module
func (msg MsgCreateReport) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateReport) Type() string { return "create_report" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateReport) ValidateBasic() error {
	if msg.Reporter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Reporter.String())
	} else if len(msg.RequestID) == 0 || msg.Validator.Empty() {
		return sdkerrors.Wrap(ErrMsgInvalid, "Request ID / validator address cannot be empty")
	} else if len(msg.DataSourceResults) == 0 || len(msg.TestCaseResults) == 0 || len(msg.AggregatedResult) == 0 {
		return sdkerrors.Wrap(ErrMsgInvalid, "lengths of the data source and test case must be greater than zero, and there must be an aggregated result")
	} else {
		return checkFees(msg.Fees.String())
	}
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateReport) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateReport) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Reporter}
}

// MsgAddReporter is a message for adding a new reporter for a validator.
type MsgAddReporter struct {
	// Validator is the validator that wishes to add a new reporter. This is the signer.
	Validator sdk.ValAddress `json:"validator"`
	// Reporter is the address to be added as a reporter to the validator.
	Reporter sdk.AccAddress `json:"reporter"`
	// Adder is the address responsible for adding the reporter
	Adder sdk.AccAddress `json:"adder"`
}

// NewMsgAddReporter is a constructor function for MsgAddReporter
func NewMsgAddReporter(validator sdk.ValAddress, reporter sdk.AccAddress, adder sdk.AccAddress) MsgAddReporter {
	return MsgAddReporter{
		Adder:     adder,
		Validator: validator,
		Reporter:  reporter,
	}
}

// Route should return the name of the module
func (msg MsgAddReporter) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddReporter) Type() string { return "add_reporter" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddReporter) ValidateBasic() error {
	if msg.Validator.Empty() || msg.Adder.Empty() || msg.Reporter.Empty() {
		return sdkerrors.Wrap(ErrReporterMsgInvalid, "The message attibutes cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Adder}
}

// MsgRemoveReporter is a message for remove an existing reporter for a validator.
type MsgRemoveReporter struct {
	// Validator is the validator that wishes to add a new reporter. This is the signer.
	Validator sdk.ValAddress `json:"validator"`
	// Reporter is the address to be added as a reporter to the validator.
	Reporter sdk.AccAddress `json:"reporter"`
	// Adder is the address responsible for adding the reporter
	Remover sdk.AccAddress `json:"remover"`
}

// NewMsgRemoveReporter is a constructor function for MsgRemoveReporter
func NewMsgRemoveReporter(validator sdk.ValAddress, reporter sdk.AccAddress, remover sdk.AccAddress) MsgRemoveReporter {
	return MsgRemoveReporter{
		Remover:   remover,
		Validator: validator,
		Reporter:  reporter,
	}
}

// Route should return the name of the module
func (msg MsgRemoveReporter) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRemoveReporter) Type() string { return "remove_reporter" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRemoveReporter) ValidateBasic() error {
	if msg.Validator.Empty() || msg.Remover.Empty() || msg.Reporter.Empty() {
		return sdkerrors.Wrap(ErrReporterMsgInvalid, "The message attibutes cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRemoveReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRemoveReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Remover}
}

// MsgCreateStrategy is a message for the information of a strategy for a yAI flow
type MsgCreateStrategy struct {
	StratID        uint64         `json:"strategy_id"`
	StratName      string         `json:"strategy_name"`
	StratFlow      []string       `json:"strategy_flow"`
	PerformanceFee uint64         `json:"performance_fee"`
	PerformanceMax uint64         `json:"performance_max"`
	WithdrawalFee  uint64         `json:"withdrawal_fee"`
	WithdrawalMax  uint64         `json:"withdrawal_max"`
	GovernanceAddr string         `json:"governance_address"`
	StrategistAddr string         `json:"strategist_address"`
	Creator        sdk.AccAddress `json:"strategy_creator"`
}

// NewMsgCreateStrategy is a constructor for the msg strategy struct
func NewMsgCreateStrategy(
	stratID uint64,
	stratName string,
	stratFlow []string,
	performanceFee uint64,
	performanceMax uint64,
	withdrawalFee uint64,
	withdrawalMax uint64,
	governanceAddr string,
	strategistAddr string,
	creator sdk.AccAddress,
) MsgCreateStrategy {
	return MsgCreateStrategy{
		StratID:        stratID,
		StratName:      stratName,
		StratFlow:      stratFlow,
		PerformanceFee: performanceFee,
		PerformanceMax: performanceMax,
		WithdrawalFee:  withdrawalFee,
		WithdrawalMax:  withdrawalMax,
		GovernanceAddr: governanceAddr,
		StrategistAddr: strategistAddr,
		Creator:        creator,
	}
}

// Route should return the name of the module
func (msg MsgCreateStrategy) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateStrategy) Type() string { return "create_strategy" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateStrategy) ValidateBasic() error {
	if len(msg.Creator) == 0 || len(msg.GovernanceAddr) == 0 || len(msg.StrategistAddr) == 0 || len(msg.StratName) == 0 {
		return sdkerrors.Wrap(ErrMsgStrategyAttrsNotFound, "The message attibutes cannot be empty")
	}
	if msg.PerformanceFee == uint64(0) || msg.PerformanceMax == uint64(0) || msg.WithdrawalFee == uint64(0) || msg.WithdrawalMax == uint64(0) {
		return sdkerrors.Wrap(ErrMsgStrategyErrorFees, "The fees are empty or invalid")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateStrategy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateStrategy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func checkFees(fees string) error {
	_, err := sdk.ParseCoins(fees)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidFeeType, err.Error())
	}
	return nil
}
