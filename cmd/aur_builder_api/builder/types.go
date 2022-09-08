package builder

import (
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/docker"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model/config"
)

type AurBuilderService struct {
	controller *docker.ContainerController
	cfg        config.DockerConfig
}

const (
	packageSuffix = "-x86_64.pkg.tar.zst"
	packagePath   = "/pkg/"
)

func NewAurBuilderService(cfg *config.DockerConfig) (*AurBuilderService, error) {
	controller, err := docker.NewContainerController(
		&docker.RegistryData{
			UseAuth:  cfg.Auth,
			Username: cfg.Username,
			Password: cfg.Password,
			Image:    cfg.ContainerImage,
		},
	)

	if err != nil {
		return nil, err
	}

	return &AurBuilderService{cfg: *cfg, controller: controller}, nil
}
