package docker

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// RunCommand runs a command and returns a reader for the output
func RunCommand(cli *client.Client, container types.Container, command []string) (int, string, error) {
	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
	}
	ctx := context.Background()

	IDResp, err := cli.ContainerExecCreate(ctx, container.ID, config)
	if err != nil {
		return 0, "", err
	}

	resp, err := cli.ContainerExecAttach(ctx, IDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return 0, "", err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Reader)
	if err != nil {
		return 0, "", err
	}

	// Get exit code
	inspectResp, err := cli.ContainerExecInspect(ctx, IDResp.ID)
	if err != nil {
		return 0, "", err
	}

	return inspectResp.ExitCode, buf.String(), err
}
