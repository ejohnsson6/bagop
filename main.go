package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	ENABLE_LABEL    = "bagop.enable"
	DB_LABEL        = "bagop.vendor"
	NAME_LABEL      = "bagop.name"
	BACKUP_LOCATION = "/backups"
)

func containerEnabled(container types.Container) bool {
	labels := container.Labels
	isEnabled := labels[ENABLE_LABEL] != "false" && labels[ENABLE_LABEL] != "" && labels[ENABLE_LABEL] != "0"
	return isEnabled
}

func findVendor(container types.Container) (string, error) {
	labelName := strings.ToLower(container.Labels[NAME_LABEL])
	allowedVendors := []string{"mysql", "postgres", "mariadb"}
	if Contains(allowedVendors, labelName) {
		return labelName, nil
	}

	imageName := strings.Split(container.Image, ":")[0]

	if Contains(allowedVendors, imageName) {
		return imageName, nil
	}
	return "", fmt.Errorf("Label:%s, Image:%s not recognized as a supported database", labelName, imageName)
}

func dumpMysql(container types.Container) {

}
func dumpPostgres(container types.Container, cli client.Client) {
	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"pg_dump", "--username=postgres", "bookit"},
	}
	ctx := context.Background()

	IDResp, err := cli.ContainerExecCreate(ctx, container.ID, config)
	if err != nil {
		log.Panic(err)
	}

	resp, err := cli.ContainerExecAttach(ctx, IDResp.ID, types.ExecStartCheck{})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(container.Names)
	// open output file
	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := resp.Reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			panic(err)
		}
	}

}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if !containerEnabled(container) {
			continue
		}
		vendor, err := findVendor(container)
		if err != nil {
			panic(err)
		}
		switch vendor {
		case "mysql", "mariadb":
			dumpMysql(container)
		case "postgres":
			dumpPostgres(container, *cli)
		}
	}
}
