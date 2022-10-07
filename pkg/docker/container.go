package docker

import (
	"archive/tar"
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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

func (c *ContainerController) CopyFromContainer(containerId string, src, dest, pkgName string) (filePathWithName string, err error) {
	stream, _, err := c.cli.CopyFromContainer(context.Background(), containerId, src)

	if err != nil {
		return "", err
	}

	defer stream.Close()

	reader := tar.NewReader(stream)

	if _, err := reader.Next(); err != nil {
		return "", err
	}

	filePathWithName = dest + pkgName

	file, err := os.Create(filePathWithName)

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	if err != nil {
		return "", err
	}

	return filePathWithName, nil
}
