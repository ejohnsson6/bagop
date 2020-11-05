package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	ENABLE_LABEL    = "bagop.enable"
	VENDOR_LABEL    = "bagop.vendor"
	NAME_LABEL      = "bagop.name"
	BACKUP_LOCATION = "/backups"
)

func containerEnabled(container types.Container) bool {
	labels := container.Labels
	isEnabled := labels[ENABLE_LABEL] != "false" && labels[ENABLE_LABEL] != "" && labels[ENABLE_LABEL] != "0"
	return isEnabled
}

func findVendor(container types.Container) (string, error) {
	labelVendor := strings.ToLower(container.Labels[VENDOR_LABEL])
	allowedVendors := []string{"mysql", "postgres", "mariadb"}
	if Contains(allowedVendors, labelVendor) {
		return labelVendor, nil
	}

	imageVendor := strings.Split(container.Image, ":")[0]

	if Contains(allowedVendors, imageVendor) {
		return imageVendor, nil
	}
	return "", fmt.Errorf("Label:%s, Image:%s not recognized as a supported database", labelVendor, imageVendor)
}

func findName(container types.Container) string {
	name := container.Labels[NAME_LABEL]
	if name == "" {
		name = container.ID
	}
	return name
}

func findEnvVar(env []string, find string) string {
	for _, e := range env {
		split := strings.Split(e, "=")
		if split[0] == find {
			return split[1]
		}
	}
	return ""
}

func dumpMysql(container types.Container, cli client.Client) error {
	return nil
}
func dumpPostgres(container types.Container, cli client.Client) error {

	inspect, err := cli.ContainerInspect(context.Background(), container.ID)
	if err != nil {
		return err
	}
	username := findEnvVar(inspect.Config.Env, "POSTGRES_USER")
	db := findEnvVar(inspect.Config.Env, "POSTGRES_DB")

	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"pg_dump", fmt.Sprintf("--username=%s", username), db},
	}
	ctx := context.Background()

	IDResp, err := cli.ContainerExecCreate(ctx, container.ID, config)
	if err != nil {
		return err
	}

	resp, err := cli.ContainerExecAttach(ctx, IDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	name := findName(container)
	// open output file
	fo, err := os.Create(name)
	if err != nil {
		return err
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
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil

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
			err = dumpMysql(container, *cli)
		case "postgres":
			err = dumpPostgres(container, *cli)
		}
		if err != nil {
			panic(err)
		}
	}
}
