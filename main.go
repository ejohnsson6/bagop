package main

import (
	"context"
	"fmt"
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

func findVendor(container types.Container) string {
	nameLabel := strings.ToLower(container.Labels[NAME_LABEL])
	allowedVendors := []string{"mysql", "postgres"}
	if Contains(allowedVendors, nameLabel) {
		return nameLabel
	}
	return "dasd"
}

func dumpMysql(container types.Container) {

}

func dumpPostgres(container types.Container) {

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
		fmt.Printf("%s %s\n", container.ID, container.Labels)
		fmt.Println(containerEnabled(container))
	}
}
