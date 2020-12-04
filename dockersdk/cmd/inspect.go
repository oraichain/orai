package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Print the docker container inspect",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			inspectCmdRun(args[0])
		}
	},
}

func inspectCmdRun(id string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	container, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(container.State.Running)

}
