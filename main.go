package main

/*
**Go client for the Docker Engine API
**https://pkg.go.dev/github.com/docker/docker/client
**
**/

//For this tutorial, it is assumed that the container is running
//Run an apline container with fllowing flags:
//--name any custom name which is going to be used later in the main function.
//-d notify docker engine to run it in detached mode i.e, run in the background
//-it tells docker to run a container interactivaly and with tty controls
//--rm flag is optional and useful in cases in which you dont want to keep them in memory after exit them
//final command will look like this:
//docker run --name awesome_mgm -d -it --rm alpine

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func modifyContainer_api(containerID string, cmdStatement []string) (err error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	cmdStatementExecuteScript := cmdStatement
	optionsCreateExecuteScript := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmdStatementExecuteScript,
	}

	response, err := cli.ContainerExecCreate(ctx, containerID, optionsCreateExecuteScript)
	if err != nil {
		panic(err)
	}

	hijackedResponse, err := cli.ContainerExecAttach(ctx, response.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	defer hijackedResponse.Close()

	return
}

func copyFilesInContainer(containerName string, scrFilePath string, dstFilePath string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containerID, status := getContainerInfo(containerName)

	if status == "Up" {
		cmdStatement := []string{"mkdir", "-p", "script"}
		modifyContainer_api(containerID, cmdStatement)

		file, err := os.Open(scrFilePath)
		if err != nil {
			panic(err)
		}
		err = cli.CopyToContainer(context.Background(), containerID, dstFilePath, bufio.NewReader(file), types.CopyToContainerOptions{
			AllowOverwriteDirWithFile: true,
		})
		if err != nil {
			panic(err)
		}
	}
}

func getContainerInfo(containerName string) (containerID string, status string) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	filters := filters.NewArgs()
	filters.Add(
		"name", containerName,
	)

	resp, err := cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters})
	if err != nil {
		panic(err)
	}

	if len(resp) > 0 {
		containerID = resp[0].ID
		containerStatus := strings.Split(resp[0].Status, " ")
		status = containerStatus[0] //fmt.Println(status[0])
	} else {
		fmt.Printf("container '%s' does not exists\n", containerName)
	}

	return
}

func main() {
	copyFilesInContainer("awesome_mgm", "file.tar.gz", "/script")
}
