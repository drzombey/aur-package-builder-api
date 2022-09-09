package builder

import (
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
)

type AurBuilderService struct {
	controller *docker.ContainerController
}

const (
	packageSuffix = "-x86_64.pkg.tar.zst"
	packagePath   = "/pkg/"
)

func NewAurBuilderService(dc *docker.ContainerController) (*AurBuilderService, error) {

	return &AurBuilderService{controller: dc}, nil
}
