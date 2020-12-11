package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

const (
	imageName = "python:3.7-alpine"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a docker container",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			createCmdRun(args[0])
		}
	},
}

func createCmdRun(containerName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      imageName,
		WorkingDir: "/.oraifiles",
		Tty:        true,
	}, nil, nil, nil, containerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp)

	// if run then remove, need waiting
	// statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// case <-statusCh:
	// }

}
