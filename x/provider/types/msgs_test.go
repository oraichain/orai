package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgRoute(t *testing.T) {
	require.Equal(t, "provider", MsgCreateAIDataSource{}.Route())
	require.Equal(t, "provider", MsgEditAIDataSource{}.Route())
	require.Equal(t, "provider", MsgCreateOracleScript{}.Route())
	require.Equal(t, "provider", MsgEditOracleScript{}.Route())
	require.Equal(t, "provider", MsgCreateTestCase{}.Route())
	require.Equal(t, "provider", MsgEditTestCase{}.Route())
}
