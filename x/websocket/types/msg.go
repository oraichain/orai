package types

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// )

// // MsgCreateStrategy is a message for the information of a strategy for a yAI flow
// type MsgCreateStrategy struct {
// 	StratID        uint64         `json:"strategy_id"`
// 	StratName      string         `json:"strategy_name"`
// 	StratFlow      []string       `json:"strategy_flow"`
// 	PerformanceFee uint64         `json:"performance_fee"`
// 	PerformanceMax uint64         `json:"performance_max"`
// 	WithdrawalFee  uint64         `json:"withdrawal_fee"`
// 	WithdrawalMax  uint64         `json:"withdrawal_max"`
// 	GovernanceAddr string         `json:"governance_address"`
// 	StrategistAddr string         `json:"strategist_address"`
// 	Creator        sdk.AccAddress `json:"strategy_creator"`
// }

// // NewMsgCreateStrategy is a constructor for the msg strategy struct
// func NewMsgCreateStrategy(
// 	stratID uint64,
// 	stratName string,
// 	stratFlow []string,
// 	performanceFee uint64,
// 	performanceMax uint64,
// 	withdrawalFee uint64,
// 	withdrawalMax uint64,
// 	governanceAddr string,
// 	strategistAddr string,
// 	creator sdk.AccAddress,
// ) MsgCreateStrategy {
// 	return MsgCreateStrategy{
// 		StratID:        stratID,
// 		StratName:      stratName,
// 		StratFlow:      stratFlow,
// 		PerformanceFee: performanceFee,
// 		PerformanceMax: performanceMax,
// 		WithdrawalFee:  withdrawalFee,
// 		WithdrawalMax:  withdrawalMax,
// 		GovernanceAddr: governanceAddr,
// 		StrategistAddr: strategistAddr,
// 		Creator:        creator,
// 	}
// }

// // Route should return the name of the module
// func (msg MsgCreateStrategy) Route() string { return RouterKey }

// // Type should return the action
// func (msg MsgCreateStrategy) Type() string { return "create_strategy" }

// // ValidateBasic runs stateless checks on the message
// func (msg MsgCreateStrategy) ValidateBasic() error {
// 	if len(msg.Creator) == 0 || len(msg.GovernanceAddr) == 0 || len(msg.StrategistAddr) == 0 || len(msg.StratName) == 0 {
// 		return sdkerrors.Wrap(ErrMsgStrategyAttrsNotFound, "The message attibutes cannot be empty")
// 	}
// 	if msg.PerformanceFee == uint64(0) || msg.PerformanceMax == uint64(0) || msg.WithdrawalFee == uint64(0) || msg.WithdrawalMax == uint64(0) {
// 		return sdkerrors.Wrap(ErrMsgStrategyErrorFees, "The fees are empty or invalid")
// 	}
// 	return nil
// }

// // GetSignBytes encodes the message for signing
// func (msg MsgCreateStrategy) GetSignBytes() []byte {
// 	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
// }

// // GetSigners defines whose signature is required
// func (msg MsgCreateStrategy) GetSigners() []sdk.AccAddress {
// 	return []sdk.AccAddress{msg.Creator}
// }
