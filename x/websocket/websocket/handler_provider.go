package websocket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider"
)

// verifyProviderMessageType checks if the message type from the system is from the provider module or not
func verifyProviderMessageType(messageType string, log sdk.ABCIMessageLog) (string, bool) {
	switch messageType {
	case (provider.MsgCreateAIDataSource{}).Type():

		// collect data source name to copy the file onto the container
		dSource, err := GetEventValue(log, provider.EventTypeSetDataSource, provider.AttributeDataSourceName)
		if err != nil {
			return "", false
		}
		return getDSourcePath(dSource), true

	case (provider.MsgCreateOracleScript{}).Type():
		// collect oscript name to copy the file onto the container
		oScript, err := GetEventValue(log, provider.EventTypeSetOracleScript, provider.AttributeOracleScriptName)
		if err != nil {
			return "", false
		}
		return getOScriptPath(oScript), true

	case (provider.MsgCreateTestCase{}).Type():
		// collect test case name to copy the file onto the container
		tCase, err := GetEventValue(log, provider.EventTypeCreateTestCase, provider.AttributeTestCaseName)
		if err != nil {
			return "", false
		}
		return getTCasePath(tCase), true

	case (provider.MsgEditAIDataSource{}).Type():
		// collect data source name to copy the file onto the container
		dSource, err := GetEventValue(log, provider.EventTypeEditDataSource, provider.AttributeDataSourceName)
		if err != nil {
			return "", false
		}
		return getDSourcePath(dSource), true

	case (provider.MsgEditOracleScript{}).Type():
		// collect oscript name to copy the file onto the container
		oScript, err := GetEventValue(log, provider.EventTypeEditOracleScript, provider.AttributeOracleScriptName)
		if err != nil {
			return "", false
		}
		return getOScriptPath(oScript), true

	case (provider.MsgEditTestCase{}).Type():
		// collect test case name to copy the file onto the container
		tCase, err := GetEventValue(log, provider.EventTypeEditTestCase, provider.AttributeTestCaseName)
		if err != nil {
			return "", false
		}
		return getTCasePath(tCase), true

	default:
		return "", false
	}
}
