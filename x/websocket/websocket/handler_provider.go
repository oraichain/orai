package websocket

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/docker/docker/client"
)

func handleProviderMsgLog(c *Context, l *Logger, log sdk.ABCIMessageLog) {
	l.Info(":delivery_truck: Processing incoming provider event")

	// create a new go routine to handle the provider message
	go func(l *Logger) {
		// create new context for the python container
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			l.Error(":skull: Failed to create new context and client for the python container: %s", err.Error())
		}

		// install new requirements for the given new script
		InstallRequirements(ctx, cli, []string{"pip", "install", "-r", "requirements.txt"})
	}(l.With())
}
