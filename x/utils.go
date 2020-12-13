package x

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const (
	fileDir = "/.oraifiles"
)

func ExecPythonFile(id string, file string, input []string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	if !CheckExistsContainer(cli, "python") {
		//create container
		workDir, err := os.Getwd()
		if err != nil {
			return "", err
		}

		err = CreateContainer(cli)
		if err != nil {
			return "", err
		}

		pythonDir := path.Join(workDir, fileDir)
		fmt.Println("python dir: ", pythonDir)
		files, err := ioutil.ReadDir(pythonDir)
		if err != nil {
			return "", err
		}
		for _, f := range files {
			if !f.IsDir() {
				fmt.Println("file copied: ", path.Join(pythonDir, f.Name()))
				CopyFileToContainer("python", path.Join(pythonDir, f.Name()))
			}
		}
	}

	fileName := filepath.Base(file)
	CopyFileToContainer("python", file)

	resp, err := cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          append([]string{"python", fileName}, input...),
	})
	if err != nil {
		return "", err
	}

	logResp, err := cli.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}

	var buf, error bytes.Buffer
	stdcopy.StdCopy(&buf, &error, logResp.Reader)
	return string(buf.Bytes()), nil
}

func CopyFileToContainer(id string, filePath string) {
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
	cli.CopyToContainer(ctx, id, fileDir, &buf, types.CopyToContainerOptions{})
}

func CheckExistsContainer(cli *client.Client, id string) bool {
	opts := types.ContainerListOptions{All: true}

	opts.Filters = filters.NewArgs()
	opts.Filters.Add("name", id)

	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, opts)
	if err != nil {
		panic(err)
	}

	return len(containers) > 0
}

func CreateContainer(cli *client.Client) error {
	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, "python:3.7-alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      "python:3.7-alpine",
		WorkingDir: fileDir,
		Tty:        true,
	}, nil, nil, nil, "python")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	//exec import requirements
	restInstall, err := cli.ContainerExecCreate(ctx, "python", types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"pip", "install", "-r", " requirements.txt"},
	})
	if err != nil {
		return err
	}
	_, err = cli.ContainerExecAttach(ctx, restInstall.ID, types.ExecStartCheck{})

	if err != nil {
		return err
	}

	return nil
}
