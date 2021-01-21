package websocket

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/docker/docker/client"
	"github.com/oraichain/orai/x/provider"
)

func handleProviderMsgLog(c *Context, l *Logger, log sdk.ABCIMessageLog, scriptPath string) {
	l.Info(":delivery_truck: Processing incoming provider event")

	// create a new go routine to handle the provider message
	go func(l *Logger) {
		// create new context for the python container
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			l.Error(":skull: Failed to create new context and client for the python container: %s", err.Error())
		}

		// verify if the python container exists or not. If not then we create new
		isExist, err := CheckExistsContainer(cli, "python")
		if err != nil {
			l.Error(":skull: Cannot check if the container exists or not: %s", err.Error())
		}
		if !isExist {
			//create container
			l.Info(":question_mark: container not exist yet")
			err = CreateContainer(ctx, cli)
			if err != nil {
				l.Error(":skull: Failed to create new python container for provider module: %s", err.Error())
			}
		}
		// install new requirements for the given new script
		CopyFileToContainer(ctx, cli, "python", scriptPath)
		err = InstallRequirements(ctx, cli, []string{"pipreqs", "--force"})
		if err != nil {
			l.Error(":skull: Failed to install requirements: %s", err.Error())
		}
		err = InstallRequirements(ctx, cli, []string{"pip", "install", "-r", "requirements.txt"})
		if err != nil {
			l.Error(":skull: Failed to install requirements: %s", err.Error())
		}
	}(l.With())
}

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
