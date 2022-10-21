package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (c *ContainerController) RemoveContainerById(containerId string) (err error) {
	cnt, err := c.ContainerById(containerId)

	if err != nil {
		return err
	}

	_, err = c.WaitForContainer(containerId, container.WaitConditionNotRunning)

	if err != nil {
		return err
	}

	err = c.cli.ContainerRemove(context.Background(), cnt.ID, types.ContainerRemoveOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (c *ContainerController) CleanContainersByImage(image string) (err error) {
	containers, err := c.ContainersByImage(image)

	if err != nil {
		return err
	}

	for _, cnt := range containers {
		_, err := c.WaitForContainer(cnt.ID, container.WaitConditionNotRunning)

		if err != nil {
			continue
		}

		err = c.cli.ContainerRemove(context.Background(), cnt.ID, types.ContainerRemoveOptions{})

		if err != nil {
			continue
		}
	}

	return nil
}
