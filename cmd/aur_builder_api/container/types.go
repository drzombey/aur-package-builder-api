package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type ContainerController struct {
	cli *client.Client
}

type VolumeMount struct {
	HostPath string
	Volume   *types.Volume
}

func NewContainerController() (c *ContainerController, err error) {
	c = new(ContainerController)

	c.cli, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		logrus.Errorf("Cannot init container client [error: %s]", err)
		return nil, err
	}

	return c, nil
}
