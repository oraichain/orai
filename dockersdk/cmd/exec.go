package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Print the docker container exec",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			execCmdRun(args[0], args[1])
		}
	},
}

func execCmdRun(id string, param string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"python", "-c", param},
	})
	if err != nil {
		panic(err)
	}

	logResp, err := cli.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, logResp.Reader)
	fmt.Println(buf.String())
}
