package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/swexbe/bagop/internal/pkg/file"
)

// RunCommand runs a command and writes the output to a file
func RunCommand(cli *client.Client, container types.Container, command []string, fileName string) (int, error) {
	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
	}
	ctx := context.Background()

	IDResp, err := cli.ContainerExecCreate(ctx, container.ID, config)
	if err != nil {
		return 0, err
	}

	resp, err := cli.ContainerExecAttach(ctx, IDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return 0, err
	}

	_, err = file.ReaderToFile(resp.Reader, fileName)
	if err != nil {
		return 0, err
	}

	// Get exit code
	inspectResp, err := cli.ContainerExecInspect(ctx, IDResp.ID)
	if err != nil {
		return 0, err
	}

	return inspectResp.ExitCode, nil
}
