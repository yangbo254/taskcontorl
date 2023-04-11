package node

import (
	"bytes"
	"context"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type DockerCmd struct {
}

func NewDockerCmd() *DockerCmd {
	return &DockerCmd{}
}

func (cmd *DockerCmd) PullImage(imageUrl string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, imageUrl, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	newStr := buf.String()
	return newStr, nil
}

func (cmd *DockerCmd) ListContainer() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (cmd *DockerCmd) ConvContainersToMap(containers []types.Container) map[string]types.Container {
	result := make(map[string]types.Container)
	for _, v := range containers {
		result[v.ID] = v
	}
	return result
}

func (cmd *DockerCmd) CreateContainer(imageUrl string, ContainerCmd []string, volumeMounts map[string]string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	hostConfig := &containertypes.HostConfig{}
	for k, v := range volumeMounts {
		mountInfo := mount.Mount{
			Type:   mount.TypeBind,
			Source: k,
			Target: v,
		}
		hostConfig.Mounts = append(hostConfig.Mounts, mountInfo)
	}

	createResp, err := cli.ContainerCreate(ctx, &containertypes.Config{
		Image: imageUrl,
		Cmd:   ContainerCmd,
	}, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
		return createResp.ID, err
	}
	return createResp.ID, nil
}

func (cmd *DockerCmd) StopContainer(containerId string, force bool) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	if force {
		if err := cli.ContainerKill(ctx, containerId, ""); err != nil {
			return err
		}
	} else {
		noWaitTimeout := 0 // to not wait for the container to exit gracefully
		if err := cli.ContainerStop(ctx, containerId, containertypes.StopOptions{Timeout: &noWaitTimeout}); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *DockerCmd) RmContainer(containerId string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	if err := cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	return nil
}
