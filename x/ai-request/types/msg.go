package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreateOracleScript defines a CreateOracleScript message
type MsgCreateOracleScript struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
}

// NewMsgCreateOracleScript is a constructor function for MsgCreateOracleScript
func NewMsgCreateOracleScript(name string, code []byte, owner sdk.AccAddress, des string) MsgCreateOracleScript {
	return MsgCreateOracleScript{
		Name:        name,
		Code:        code,
		Owner:       owner,
		Description: des,
	}
}

// Route should return the name of the module
func (msg MsgCreateOracleScript) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateOracleScript) Type() string { return "set_oscript" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateOracleScript) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Name) == 0 || len(msg.Code) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Code cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgEditOracleScript defines a message for editing a oScript in the store
type MsgEditOracleScript struct {
	OldName     string         `json:"old_name"`
	NewName     string         `json:"new_name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
}

// NewMsgEditOracleScript is a constructor function for MsgEditOracleScript
func NewMsgEditOracleScript(oldName string, newName string, code []byte, owner sdk.AccAddress, des string) MsgEditOracleScript {
	return MsgEditOracleScript{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Code:        code,
		Owner:       owner,
	}
}

// Route should return the name of the module
func (msg MsgEditOracleScript) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEditOracleScript) Type() string { return "edit_oscript" }

// ValidateBasic runs stateless checks on the message
func (msg MsgEditOracleScript) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OldName) == 0 || len(msg.Code) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Code cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgEditOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgCreateAIDataSource defines a MsgCreateAIDataSource message
type MsgCreateAIDataSource struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        string         `json:"transaction_fee"`
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

// Route should return the name of the module
func (msg MsgCreateAIDataSource) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateAIDataSource) Type() string { return "set_datasource" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateAIDataSource) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Name) == 0 || len(msg.Code) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Code cannot be empty")
	}
	return checkFees(msg.Fees)
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateAIDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateAIDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgEditAIDataSource defines a message for editing a data source in the store
type MsgEditAIDataSource struct {
	OldName     string         `json:"old_name"`
	NewName     string         `json:"new_name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        string         `json:"new_transaction_fee"`
}

// NewMsgEditAIDataSource is a constructor function for MsgEditAIDataSource
func NewMsgEditAIDataSource(oldName string, newName string, code []byte, owner sdk.AccAddress, fees string, des string) MsgEditAIDataSource {
	return MsgEditAIDataSource{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Code:        code,
		Owner:       owner,
		Fees:        fees,
	}
}

// Route should return the name of the module
func (msg MsgEditAIDataSource) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEditAIDataSource) Type() string { return "edit_datasource" }

// ValidateBasic runs stateless checks on the message
func (msg MsgEditAIDataSource) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OldName) == 0 || len(msg.Code) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Code cannot be empty")
	}
	return checkFees(msg.Fees)
}

// GetSignBytes encodes the message for signing
func (msg MsgEditAIDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditAIDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgSetKYCRequest defines message for a KYC request
type MsgSetKYCRequest struct {
	ImageHash    string          `json:"image_hash"`
	ImageName    string          `json:"image_name"`
	MsgAIRequest MsgSetAIRequest `json:"msg_set_ai_request"`
}

// NewMsgSetKYCRequest is a constructor function for MsgSetKYCRequest
func NewMsgSetKYCRequest(imageHash string, imageName string, msgSetAIRequest MsgSetAIRequest) MsgSetKYCRequest {
	return MsgSetKYCRequest{
		ImageHash:    imageHash,
		ImageName:    imageName,
		MsgAIRequest: msgSetAIRequest,
	}
}

// Route should return the name of the module
func (msg MsgSetKYCRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetKYCRequest) Type() string { return "set_kyc_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetKYCRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	err := msg.MsgAIRequest.ValidateBasic()
	if err != nil {
		return err
	}
	if len(msg.ImageHash) == 0 || len(msg.ImageName) == 0 {
		return sdkerrors.Wrap(ErrImageFailedToUnzip, "Image name / hash is not valid")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetKYCRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetKYCRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}

// MsgSetAIRequest defines message for an AI request
type MsgSetAIRequest struct {
	RequestID        string         `json:"request_id"`
	OracleScriptName string         `json:"oscript_name"`
	Creator          sdk.AccAddress `json:"creator"`
	ValidatorCount   int            `json:"validator_count"`
	Fees             string         `json:"transaction_fee"`
	Input            string         `json:"request_input"`
	ExpectedOutput   string         `json:"expected_output"`
}

// NewMsgSetAIRequest is a constructor function for NewMsgSetAIRequest
func NewMsgSetAIRequest(requestID string, oscriptName string, creator sdk.AccAddress, fees string, valCount int, input string, expectedOutput string) MsgSetAIRequest {
	return MsgSetAIRequest{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		ValidatorCount:   valCount,
		Fees:             fees,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}

// Route should return the name of the module
func (msg MsgSetAIRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetAIRequest) Type() string { return "set_ai_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetAIRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OracleScriptName) == 0 || msg.ValidatorCount <= 0 {
		return sdkerrors.Wrap(ErrNameIsEmpty, "Name or / and validator count cannot be empty")
	}
	return checkFees(msg.Fees)
}

// GetSignBytes encodes the message for signing
func (msg MsgSetAIRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetAIRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgSetPriceRequest defines message for a price prediction request
type MsgSetPriceRequest struct {
	MsgAIRequest MsgSetAIRequest `json:"msg_set_ai_request"`
}

// NewMsgSetPriceRequest is a constructor function for MsgSetPriceRequest
func NewMsgSetPriceRequest(msgSetAIRequest MsgSetAIRequest) MsgSetPriceRequest {
	return MsgSetPriceRequest{
		MsgAIRequest: msgSetAIRequest,
	}
}

// Route should return the name of the module
func (msg MsgSetPriceRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetPriceRequest) Type() string { return "set_price_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetPriceRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	err := msg.MsgAIRequest.ValidateBasic()
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetPriceRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetPriceRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}

// MsgCreateTestCase defines message for an AI request test case
type MsgCreateTestCase struct {
	Name        string         `json:"test_case_name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        []byte         `json:"code"`
	Fees        string         `json:"transaction_fee"`
	Description string         `json:"description"`
}

// NewMsgCreateTestCase is a constructor function for MsgCreateTestCase
func NewMsgCreateTestCase(name string, code []byte, owner sdk.AccAddress, fees string, des string) MsgCreateTestCase {
	return MsgCreateTestCase{
		Name:        name,
		Description: des,
		Owner:       owner,
		Code:        code,
		Fees:        fees,
	}
}

// Route should return the name of the module
func (msg MsgCreateTestCase) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateTestCase) Type() string { return "set_test_case" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateTestCase) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Code) == 0 {
		return sdkerrors.Wrap(ErrNameIsEmpty, "Name or/and code cannot be empty")
	}
	return checkFees(msg.Fees)
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateTestCase) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateTestCase) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgEditTestCase defines a message for editing a test case in the store
type MsgEditTestCase struct {
	OldName     string         `json:"old_name"`
	NewName     string         `json:"new_name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        string         `json:"new_transaction_fee"`
}

// NewMsgEditTestCase is a constructor function for MsgEditTestCase
func NewMsgEditTestCase(oldName string, newName string, code []byte, owner sdk.AccAddress, fees string, des string) MsgEditTestCase {
	return MsgEditTestCase{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Code:        code,
		Owner:       owner,
		Fees:        fees,
	}
}

// Route should return the name of the module
func (msg MsgEditTestCase) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEditTestCase) Type() string { return "edit_test_case" }

// ValidateBasic runs stateless checks on the message
func (msg MsgEditTestCase) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OldName) == 0 || len(msg.Code) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Code cannot be empty")
	}
	return checkFees(msg.Fees)
}

// GetSignBytes encodes the message for signing
func (msg MsgEditTestCase) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditTestCase) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
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
