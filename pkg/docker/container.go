package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"io"
)

func (c *ContainerController) WaitForContainer(containerId string, condition container.WaitCondition) (state int64, err error) {
	resultC, errC := c.cli.ContainerWait(context.Background(), containerId, condition)
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

	containers, err = c.cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: args,
	})

	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (c *ContainerController) ContainerById(containerId string) (cont *types.Container, err error) {
	args := filters.NewArgs()
	args.Add("id", containerId)

	containers, err := c.cli.ContainerList(context.Background(), container.ListOptions{
		Limit:   1,
		Filters: args,
	})

	if err != nil {
		return nil, err
	}

	cont = &containers[0]
	return cont, nil
}

func (c *ContainerController) CopyFromContainer(containerId string, src string) (stream io.ReadCloser, err error) {
	stream, _, err = c.cli.CopyFromContainer(context.Background(), containerId, src)

	if err != nil {
		return nil, err
	}

	return stream, nil
}
