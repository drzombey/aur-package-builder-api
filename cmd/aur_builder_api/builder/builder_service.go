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
			Url:      s.cfg.RegistryUrl,
		},
	)

	status = "error"

	if err != nil {
		return status, err
	}

	err = controller.EnsureImage("gitlab.powerofcloud.de:5050/tim/aur-package-build:latest")

	if err != nil {
		return status, err
	}

	containerId, err := controller.RunContainer(
		"gitlab.powerofcloud.de:5050/tim/aur-package-build:latest",
		[]string{"sh", "-c", "chmod +x /aur.sh && ./aur.sh google-chrome-dev"},
		nil,
	)

	if err != nil {
		return "err", err
	}

	status = fmt.Sprintf("started container with id: %s", containerId)
	return status, nil
}
