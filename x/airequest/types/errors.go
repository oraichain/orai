package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrRequestNotFound        = sdkerrors.Register(ModuleName, 1, "The request is not found")
	ErrEditorNotAuthorized    = sdkerrors.Register(ModuleName, 2, "Unauthorized ! This behaviour is only allowed for the owner")
	ErrNotEnoughValidators    = sdkerrors.Register(ModuleName, 3, "Not enough validators to execute the request")
	ErrCannotRandomValidators = sdkerrors.Register(ModuleName, 4, "Cannot random a list of validators")
	ErrCannotFindValidator    = sdkerrors.Register(ModuleName, 5, "Cannot find the correct validator")
	ErrSeedinitiation         = sdkerrors.Register(ModuleName, 6, "The seeding initiation process has an error")
	ErrValidatorsHaveNoVotes  = sdkerrors.Register(ModuleName, 7, "Total voting power of all validators is zero")
	ErrFailedToModulo         = sdkerrors.Register(ModuleName, 8, "Tbere is an error while sampling the validators from modulo")
	ErrNeedMoreFees           = sdkerrors.Register(ModuleName, 9, "Total fee is higher than the fee given")
	ErrRequestInvalid         = sdkerrors.Register(ModuleName, 10, "The transaction fee is invalid")
	ErrRequestFeesInvalid     = sdkerrors.Register(ModuleName, 11, "The request fee is invalid")
	ErrOScriptNotFound        = sdkerrors.Register(ModuleName, 12, "The oracle script is not found")
)
