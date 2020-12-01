package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/swexbe/bagop/internal/pkg/utility"
)

// GetEnabledContainers returns all containers which have the enabled label
func GetEnabledContainers(cli *client.Client) ([]types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	// Filter out only the enabled containers
	filteredContainers := []types.Container{}
	for _, container := range containers {
		if containerEnabled(container) {
			filteredContainers = append(filteredContainers, container)
		}
	}

	return filteredContainers, nil
}

// GetEnv returns the environment variables of a docker container
func GetEnv(cli *client.Client, container types.Container) []string {
	inspect, _ := cli.ContainerInspect(context.Background(), container.ID)
	return inspect.Config.Env
}

func containerEnabled(container types.Container) bool {
	labels := container.Labels
	isEnabled := labels[enableLabel] != "false" && labels[enableLabel] != "" && labels[enableLabel] != "0"
	return isEnabled
}

// FindVendor returns the DB vendor of a given container
// Returns error if no supported vendor could be identified
func FindVendor(container types.Container) (string, error) {
	labelVendor := strings.ToLower(container.Labels[vendorLabel])
	allowedVendors := []string{"mysql", "postgres", "mariadb"}
	if utility.Contains(allowedVendors, labelVendor) {
		return labelVendor, nil
	}

	imageVendor := strings.Split(container.Image, ":")[0]

	if utility.Contains(allowedVendors, imageVendor) {
		return imageVendor, nil
	}
	return "", fmt.Errorf("Label:%s, Image:%s not recognized as a supported database", labelVendor, imageVendor)
}

// FindName returns the name set by the name label, if no name label exists instead returns the docker ID
func FindName(container types.Container) string {
	name := container.Labels[nameLabel]
	if name == "" {
		name = container.ID
	}
	return name
}
