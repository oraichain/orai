package websocket

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

	isExist, err := CheckExistsContainer(cli, "python")
	if err != nil {
		return "", err
	}

	if !isExist {
		//create container
		err = CreateContainer(ctx, cli)
		if err != nil {
			return "", err
		}
	}

	fileName := filepath.Base(file)
	CopyFileToContainer(ctx, cli, "python", file)

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

func CopyFileToContainer(ctx context.Context, cli *client.Client, id string, filePath string) error {
	fmt.Println("ready to copy file to container: ", filePath)
	fileName := filepath.Base(filePath)
	content, err := ioutil.ReadFile(filePath)
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: fileName,            // filename
		Mode: 0777,                // permissions
		Size: int64(len(content)), // filesize
	})
	if err != nil {
		return err
	}
	tw.Write([]byte(content))
	tw.Close()

	// use &buf as argument for content in CopyToContainer
	err = cli.CopyToContainer(ctx, id, fileDir, &buf, types.CopyToContainerOptions{})
	return err
}

func CheckExistsContainer(cli *client.Client, id string) (bool, error) {
	opts := types.ContainerListOptions{All: true}

	opts.Filters = filters.NewArgs()
	opts.Filters.Add("name", id)

	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, opts)
	if err != nil {
		return false, err
	}

	return len(containers) > 0, nil
}

func CreateContainer(ctx context.Context, cli *client.Client) error {
	reader, err := cli.ImagePull(ctx, "python:3.7-alpine", types.ImagePullOptions{})
	if err != nil {
		return err
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

	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	pythonDir := path.Join(workDir, fileDir)
	fmt.Println("python dir: ", pythonDir)

	// check if the directory exists or not.
	if _, err := os.Stat(pythonDir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist, we will create new with full permission 777
			err = os.MkdirAll(pythonDir, 777)
			if err != nil {
				return err
			}
		} else {
			// other error like permissions, etc, then we will not handle this case
			return err
		}
	}

	files, err := ioutil.ReadDir(pythonDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {
			CopyFileToContainer(ctx, cli, "python", path.Join(pythonDir, f.Name()))
		}
	}
	fmt.Println("ready to install requirements ...")
	//exec import requirements

	// install requirements for the python container

	// install g++
	if err = InstallRequirements(ctx, cli, []string{"apk", "add", "g++"}); err != nil {
		return err
	}

	// install pipreqs
	if err = InstallRequirements(ctx, cli, []string{"pip", "install", "pipreqs"}); err != nil {
		return err
	}

	// run pipreqs
	if err = InstallRequirements(ctx, cli, []string{"pipreqs", "--force"}); err != nil {
		return err
	}

	// install requirements.txt
	if err = InstallRequirements(ctx, cli, []string{"pip", "install", "-r", "requirements.txt"}); err != nil {
		return err
	}

	return nil
}

// InstallRequirements installs all the modules in the requirements.txt file
func InstallRequirements(ctx context.Context, cli *client.Client, cmd []string) error {
	restInstall, err := cli.ContainerExecCreate(ctx, "python", types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		//Cmd:          []string{"pip", "freeze", "-r", "requirements.txt", "&&", "pip", "install", "-r", "requirements.txt"},
		Cmd: cmd,
	})
	if err != nil {
		return err
	}
	logResp, err := cli.ContainerExecAttach(ctx, restInstall.ID, types.ExecStartCheck{})

	if err != nil {
		return err
	}

	var buf, error bytes.Buffer
	stdcopy.StdCopy(&buf, &error, logResp.Reader)
	return nil
}
