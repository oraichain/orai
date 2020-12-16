package cmd

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(existsCmd)
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Print the docker container exec",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			execCmdRun(args[0], args[1], []string{args[2]})
		} else if len(args) > 1 {
			execCmdRun(args[0], args[1], []string{})
		}
	},
}

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Print the docker container exec",
	Run: func(cmd *cobra.Command, args []string) {
		copyCmdRun(args[0], args[1])
	},
}

var existsCmd = &cobra.Command{
	Use:   "exists",
	Short: "Print the docker container exec",
	Run: func(cmd *cobra.Command, args []string) {
		CheckExistsContainer(args[0])
	},
}

func copyCmdRun(id string, filePath string) {
	ctx := context.Background()
	fileName := filepath.Base(filePath)
	content, err := ioutil.ReadFile(filePath)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: fileName,            // filename
		Mode: 0777,                // permissions
		Size: int64(len(content)), // filesize
	})
	if err != nil {
		panic(err)
	}
	tw.Write([]byte(content))
	tw.Close()

	// use &buf as argument for content in CopyToContainer
	cli.CopyToContainer(ctx, id, "/.oraifiles", &buf, types.CopyToContainerOptions{})
}

func CheckExistsContainer(id string) {
	opts := types.ContainerListOptions{All: true}

	opts.Filters = filters.NewArgs()
	opts.Filters.Add("name", id)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(containers) > 0)
}

func execCmdRun(id string, filename string, input []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		//Cmd:          append([]string{"python", "coingecko_eth.py"}, input...),
		Cmd: []string{"pip", "install", "-r", "requirements.txt"},
	})
	if err != nil {
		panic(err)
	}

	logResp, err := cli.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	//buf := &bytes.Buffer
	var out bytes.Buffer
	var error bytes.Buffer
	stdcopy.StdCopy(&out, &error, logResp.Reader)
	fmt.Println(string(out.Bytes()))
}
