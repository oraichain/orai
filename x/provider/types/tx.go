package types

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// regex allow only alphabet, numeric and underscore characters
var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_ ]*$`).MatchString

func checkFees(fees string) error {
	_, err := sdk.ParseCoinsNormalized(fees)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidFeeType, err.Error())
	}
	return nil
}
