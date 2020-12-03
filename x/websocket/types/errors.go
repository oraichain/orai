package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrCannotFindValidator = sdkerrors.Register(ModuleName, 1, "Cannot find the correct validator")
	ErrInvalidFeeType      = sdkerrors.Register(ModuleName, 3, "The transaction fee is invalid")
	ErrReporterExists      = sdkerrors.Register(ModuleName, 4, "The reporter of this validator already exists")
	ErrReporterNotFound    = sdkerrors.Register(ModuleName, 5, "The reporter of this validator cannot be found")
	ErrReporterMsgInvalid  = sdkerrors.Register(ModuleName, 6, "The reporter of this validator cannot be found")
	ErrMsgReportInvalid    = sdkerrors.Register(ModuleName, 7, "The msg create report is invalid")
)
