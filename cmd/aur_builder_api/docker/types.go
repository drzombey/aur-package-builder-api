package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ContainerController struct {
	cli          *client.Client
	RegistryData *RegistryData
}

type VolumeMount struct {
	HostPath string
	Volume   *types.Volume
}

type RegistryData struct {
	Url      string
	Username string
	Password string
	UseAuth  bool
}

func NewContainerController(registry *RegistryData) (c *ContainerController, err error) {
	c = &ContainerController{RegistryData: registry}

	c.cli, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}

	return c, nil
}
