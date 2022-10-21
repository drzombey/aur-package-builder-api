package docker

import (
	"errors"

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
	Image    string
	Username string
	Password string
	UseAuth  bool
}

func NewContainerController(registry *RegistryData) (c *ContainerController, err error) {

	if registry.Image == "" {
		return nil, errors.New("container controller cannot be initialized cause image name not found")
	}

	c = &ContainerController{RegistryData: registry}

	c.cli, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}

	return c, nil
}
