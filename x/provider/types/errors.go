package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrOracleScriptNotFound     = sdkerrors.Register(ModuleName, 1, "The oScript does not exist")
	ErrDataSourceNotFound       = sdkerrors.Register(ModuleName, 2, "The data source does not exist")
	ErrResultDoesNotExist       = sdkerrors.Register(ModuleName, 3, "The result does not exist")
	ErrOracleScriptNameExists   = sdkerrors.Register(ModuleName, 4, "The oScript name already exists")
	ErrDataSourceNameExists     = sdkerrors.Register(ModuleName, 5, "The data source name already exists")
	ErrCodeNotFound             = sdkerrors.Register(ModuleName, 6, "The code file does not exist")
	ErrNameIsEmpty              = sdkerrors.Register(ModuleName, 7, "The oScript name cannot be empty")
	ErrRequestNotFound          = sdkerrors.Register(ModuleName, 8, "The request is not found")
	ErrResultNotFound           = sdkerrors.Register(ModuleName, 9, "The result is not found")
	ErrEditorNotAuthorized      = sdkerrors.Register(ModuleName, 10, "Unauthorized ! This behaviour is only allowed for the owner")
	ErrImageFailedToUnzip       = sdkerrors.Register(ModuleName, 11, "Cannot unzip the image")
	ErrFailedToOpenFile         = sdkerrors.Register(ModuleName, 12, "Failed to open script file. The file seems to be missing or not found. Please check if the data source or test case exists")
	ErrFailedToRequestFile      = sdkerrors.Register(ModuleName, 13, "Failed to request file from IPFS")
	ErrFailedToWriteFile        = sdkerrors.Register(ModuleName, 14, "Failed to write the file from the request result")
	ErrNotEnoughValidators      = sdkerrors.Register(ModuleName, 15, "Not enough validators to execute the request")
	ErrCannotRandomValidators   = sdkerrors.Register(ModuleName, 16, "Cannot random a list of validators")
	ErrCannotFindValidator      = sdkerrors.Register(ModuleName, 17, "Cannot find the correct validator")
	ErrValidatorAlreadyReported = sdkerrors.Register(ModuleName, 18, "The validator is already reported")
	ErrTestCaseNotFound         = sdkerrors.Register(ModuleName, 19, "The test case is not found")
	ErrInvalidFeeType           = sdkerrors.Register(ModuleName, 20, "The transaction fee is invalid")
	ErrTestCaseNameExists       = sdkerrors.Register(ModuleName, 21, "The test case name already exists")
	ErrNeedMoreFees             = sdkerrors.Register(ModuleName, 22, "Total fee is higher than the fee given")
	ErrReporterExists           = sdkerrors.Register(ModuleName, 23, "The reporter of this validator already exists")
	ErrReporterNotFound         = sdkerrors.Register(ModuleName, 24, "The reporter of this validator cannot be found")
	ErrReporterMsgInvalid       = sdkerrors.Register(ModuleName, 25, "The reporter of this validator cannot be found")
	ErrMsgInvalid               = sdkerrors.Register(ModuleName, 26, "The msg is invalid")
	ErrMsgStrategyErrorFees     = sdkerrors.Register(ModuleName, 27, "The fees are invalid")
	ErrMsgStrategyAttrsNotFound = sdkerrors.Register(ModuleName, 28, "strategy attributes not found")
	ErrPaginationInputInvalid   = sdkerrors.Register(ModuleName, 29, "The page and limit inputs are invalid")
	ErrSeedinitiation           = sdkerrors.Register(ModuleName, 30, "The seeding initiation process has an error")
	ErrValidatorsHaveNoVotes    = sdkerrors.Register(ModuleName, 31, "Total voting power of all validators is zero")
	ErrFailedToModulo           = sdkerrors.Register(ModuleName, 32, "Tbere is an error while sampling the validators from modulo")
	ErrCannotGetMinimumFees     = sdkerrors.Register(ModuleName, 33, "Cannot retrieve minimum fees")
)
