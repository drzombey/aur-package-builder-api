package builder

import (
	"fmt"

	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/docker"
	"github.com/drzombey/aur-rpc-client-go/types"
)

func (s *AurBuilderService) BuildAurPackage(aurpackage *types.Package) (status string, err error) {
	controller, err := docker.NewContainerController(
		&docker.RegistryData{
			UseAuth:  s.cfg.Auth,
			Username: s.cfg.Username,
			Password: s.cfg.Password,
			Image:    s.cfg.ContainerImage,
		},
	)

	status = "error"

	if err != nil {
		return status, err
	}

	err = controller.EnsureImage()

	if err != nil {
		return status, err
	}

	containerId, err := controller.RunContainer(
		controller.RegistryData.Image,
		[]string{"sh", "-c", "./aur.sh " + aurpackage.PackageBase},
		nil,
	)

	if err != nil {
		return "err", err
	}

	status = fmt.Sprintf("started container with id: %s", containerId)
	return status, nil
}
