package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunCommand(cli *client.Client, container types.Container, command []string) (io.Reader, error) {
	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
	}
	ctx := context.Background()

	IDResp, err := cli.ContainerExecCreate(ctx, container.ID, config)
	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerExecAttach(ctx, IDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}
	return resp.Reader, nil
}
