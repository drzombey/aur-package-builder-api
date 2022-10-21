package builder

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/drzombey/aur-package-builder-api/pkg/aur"
	"github.com/sirupsen/logrus"
)

func (s *AurBuilderService) StartBuildAurPkgRoutine(aurPkg *aur.Package, dest string) (containerId string, err error) {
	containerId, err = s.buildAurPackage(aurPkg)

	if err != nil {
		return "", err
	}

	go func() {
		_, err = s.copyPackageToDestination(containerId, dest, aurPkg)

		if err != nil {
			logrus.Errorf("Copying package from container %s not possible [error: %s]", containerId, err)
		}

		err = s.cleanUpBuildEnv(containerId)

		if err != nil {
			logrus.Errorf("Cleanup container %s not possible [error: %s]", containerId, err)
		}
	}()

	return containerId, nil
}

func (s *AurBuilderService) buildAurPackage(aurPkg *aur.Package) (containerId string, err error) {

	logrus.Debug("Ensure image is available on system")
	err = s.controller.EnsureImage()

	if err != nil {
		return "", err
	}

	image := s.controller.RegistryData.Image
	command := []string{"sh", "-c", "./aur.sh " + aurPkg.PackageBase}

	logrus.Debugf("Try to start container [Image: %s, Cmd: %v]", image, command)
	containerId, err = s.controller.RunContainer(
		image,
		command,
		nil,
	)
	logrus.Debugf("Container started with id: %s", containerId)

	if err != nil {
		return "", err
	}

	return containerId, nil
}

func (s *AurBuilderService) copyPackageToDestination(containerId, dest string, aurPkg *aur.Package) (pkgPath string, err error) {
	pkgName := fmt.Sprintf("%s-%s%s", aurPkg.PackageBase, aurPkg.Version, packageSuffix)

	_, err = s.controller.WaitForContainer(containerId, container.WaitConditionNextExit)
	if err != nil {
		return "", err
	}

	pkgPath, err = s.controller.CopyFromContainer(containerId, packagePath+pkgName, dest, pkgName)

	if err != nil {
		return "", err
	}

	return pkgPath, nil
}

func (s *AurBuilderService) cleanUpBuildEnv(containerId string) (err error) {
	err = s.controller.RemoveContainerById(containerId)

	if err != nil {
		return err
	}

	return nil
}
