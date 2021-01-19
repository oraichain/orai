package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrOracleScriptNotFound   = sdkerrors.Register(ModuleName, 1, "The oScript does not exist")
	ErrDataSourceNotFound     = sdkerrors.Register(ModuleName, 2, "The data source does not exist")
	ErrOracleScriptNameExists = sdkerrors.Register(ModuleName, 4, "The oScript name already exists")
	ErrDataSourceNameExists   = sdkerrors.Register(ModuleName, 5, "The data source name already exists")
	ErrCodeNotFound           = sdkerrors.Register(ModuleName, 6, "The code file does not exist")
	ErrEmpty                  = sdkerrors.Register(ModuleName, 7, "The script or code cannot be empty")
	ErrEditorNotAuthorized    = sdkerrors.Register(ModuleName, 8, "Unauthorized ! This behaviour is only allowed for the owner")
	ErrFailedToOpenFile       = sdkerrors.Register(ModuleName, 9, "Failed to open script file. The file seems to be missing or not found. Please check if the data source or test case exists")
	ErrTestCaseNotFound       = sdkerrors.Register(ModuleName, 10, "The test case is not found")
	ErrInvalidFeeType         = sdkerrors.Register(ModuleName, 11, "The transaction fee is invalid")
	ErrTestCaseNameExists     = sdkerrors.Register(ModuleName, 12, "The test case name already exists")
	ErrPaginationInputInvalid = sdkerrors.Register(ModuleName, 13, "The page and limit inputs are invalid")
	ErrCannotGetMinimumFees   = sdkerrors.Register(ModuleName, 14, "Cannot retrieve minimum fees")
)
