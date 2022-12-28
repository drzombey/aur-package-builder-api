package builder

import (
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
)

type AurBuilderService struct {
	controller      *docker.ContainerController
	storageProvider storage.Provider
}

const (
	packageSuffix = "-x86_64.pkg.tar.zst"
	packagePath   = "/pkg/"
)

func NewAurBuilderService(dc *docker.ContainerController, sp storage.Provider) (*AurBuilderService, error) {
	return &AurBuilderService{
		controller:      dc,
		storageProvider: sp,
	}, nil
}
