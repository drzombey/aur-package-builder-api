package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func (c *ContainerController) RunContainer(image string, command []string, volumes []VolumeMount) (containerId string, err error) {
	hc := container.HostConfig{}
	var volumeMounts []mount.Mount

	for _, volume := range volumes {
		mount := mount.Mount{
			Type:   mount.TypeVolume,
			Source: volume.Volume.Name,
			Target: volume.HostPath,
		}
		volumeMounts = append(volumeMounts, mount)
	}

	hc.Mounts = volumeMounts

	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Tty:   true,
		Image: image,
		Cmd:   command,
	}, &hc, nil, nil, "")

	if err != nil {
		return "", err
	}

	err = c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
