package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

func (c *ContainerController) CleanContainersByImage(image string) (err error) {
	containers, err := c.ContainersByImage(image)

	if err != nil {
		fmt.Printf("Cannot get containers: %s", err)
		return err
	}

	for _, container := range containers {
		_, err := c.WaitForContainer(container.ID)

		if err != nil {
			fmt.Printf("Unable to wait for container %q exit: %q\n", container.ID, err)
			continue
		}

		err = c.cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})

		if err != nil {
			fmt.Printf("Unable to remove container %q: %q\n", container.ID, err)
			continue
		}
	}

	return nil
}
