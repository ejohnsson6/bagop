package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	ENABLE_LABEL    = "bagop.enable"
	VENDOR_LABEL    = "bagop.vendor"
	NAME_LABEL      = "bagop.name"
	BACKUP_LOCATION = "/backups/"
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

func readerToFile(reader io.Reader, baseDir string) error {

	os.MkdirAll(baseDir, 0755)

	timestamp := time.Now().Format(time.RFC3339)

	// open output file
	fo, err := os.Create(baseDir + "/" + timestamp)
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
		n, err := reader.Read(buf)
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

func runCommand(cli client.Client, containerID string, command []string) (io.Reader, error) {
	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
	}
	ctx := context.Background()

	IDResp, err := cli.ContainerExecCreate(ctx, containerID, config)
	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerExecAttach(ctx, IDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}
	return resp.Reader, nil
}

func dumpMysql(container types.Container, cli client.Client) error {
	inspect, err := cli.ContainerInspect(context.Background(), container.ID)
	if err != nil {
		return err
	}
	username := findEnvVar(inspect.Config.Env, "MYSQL_USER")
	password := findEnvVar(inspect.Config.Env, "MYSQL_PASSWORD")
	db := findEnvVar(inspect.Config.Env, "MYSQL_DATABASE")

	name := findName(container)

	command := []string{"mysqldump", fmt.Sprintf("-u%s", username), fmt.Sprintf("-p%s", password), "--skip-comments", "--databases", db}

	reader, err := runCommand(cli, container.ID, command)
	if err != nil {
		return err
	}
	err = readerToFile(reader, name)
	if err != nil {
		return err
	}
	return nil
}
func dumpPostgres(container types.Container, cli client.Client) error {

	inspect, err := cli.ContainerInspect(context.Background(), container.ID)
	if err != nil {
		return err
	}
	username := findEnvVar(inspect.Config.Env, "POSTGRES_USER")
	db := findEnvVar(inspect.Config.Env, "POSTGRES_DB")

	name := BACKUP_LOCATION + findName(container)

	command := []string{"pg_dump", fmt.Sprintf("--username=%s", username), db}

	reader, err := runCommand(cli, container.ID, command)
	if err != nil {
		return err
	}
	err = readerToFile(reader, name)
	if err != nil {
		return err
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
			log.Printf("Dumping MYSQL/MariaDB container %s", container.ID)
			err = dumpMysql(container, *cli)
		case "postgres":
			log.Printf("Dumping PostgreSQL container %s", container.ID)
			err = dumpPostgres(container, *cli)
		}
		if err != nil {
			panic(err)
		}
	}
}
