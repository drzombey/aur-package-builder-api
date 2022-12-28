package builder

import (
	"archive/tar"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/drzombey/aur-package-builder-api/pkg/aur"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func (s *AurBuilderService) StartBuildAurPkgRoutine(aurPkg *aur.Package, dest string) (containerId string, err error) {
	containerId, err = s.buildAurPackage(aurPkg)

	if err != nil {
		return "", err
	}

	go func() {
		err = s.savePackage(containerId, aurPkg)

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

func (s *AurBuilderService) savePackage(containerId string, aurPkg *aur.Package) (err error) {
	pkgName := fmt.Sprintf("%s-%s%s", aurPkg.PackageBase, aurPkg.Version, packageSuffix)

	_, err = s.controller.WaitForContainer(containerId, container.WaitConditionNextExit)
	if err != nil {
		return err
	}

	stream, err := s.controller.CopyFromContainer(containerId, packagePath+pkgName)

	if err != nil {
		logrus.Errorf("Got error: %s", err)
		return err
	}

	err = s.storageProvider.AddFile(context.Background(), pkgName, stream)
	if err != nil {
		logrus.Errorf("Got error: %s", err)
		return err
	}

	return nil
}

func (s *AurBuilderService) createPackageFileFromStream(stream io.ReadCloser, dest, pkgName string) (filenameWithPath string, err error) {
	reader := tar.NewReader(stream)

	if _, err := reader.Next(); err != nil {
		return "", err
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		os.MkdirAll(dest, os.ModePerm)
	}

	filenameWithPath = dest + "/" + pkgName

	file, err := os.Create(filenameWithPath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	if err != nil {
		return "", err
	}

	return filenameWithPath, nil
}

func (s *AurBuilderService) cleanUpBuildEnv(containerId string) (err error) {
	err = s.controller.RemoveContainerById(containerId)

	if err != nil {
		return err
	}

	return nil
}
