package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (c *ContainerController) WaitForContainer(containerId string) (state int64, err error) {
	resultC, errC := c.cli.ContainerWait(context.Background(), containerId, "")
	select {
	case err := <-errC:
		return 0, err
	case result := <-resultC:
		return result.StatusCode, nil
	}
}

func (c *ContainerController) ContainersByImage(image string) (containers []types.Container, err error) {
	args := filters.NewArgs()
	args.Add("ancestor", image)

	containers, err = c.cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: args,
	})

	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (c *ContainerController) ContainerById(containerId string) (container *types.Container, err error) {
	args := filters.NewArgs()
	args.Add("id", containerId)

	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{
		Limit:   1,
		Filters: args,
	})

	if err != nil {
		return nil, err
	}

	container = &containers[0]
	return container, nil
}
