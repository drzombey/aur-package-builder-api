package builder

import (
	"fmt"

	"github.com/drzombey/aur-rpc-client-go/types"
)

func (s *AurBuilderService) BuildAurPackage(aurpackage *types.Package) (containerId string, err error) {

	if err != nil {
		return "", err
	}

	err = s.controller.EnsureImage()

	if err != nil {
		return "", err
	}

	containerId, err = s.controller.RunContainer(
		s.controller.RegistryData.Image,
		[]string{"sh", "-c", "./aur.sh " + aurpackage.PackageBase},
		nil,
	)

	if err != nil {
		return "", err
	}

	return containerId, nil
}

func (s *AurBuilderService) CopyPackageToDestination(containerId, dest string, aurpkg *types.Package) (pkgPath string, err error) {
	pkgName := fmt.Sprintf("%s-%s%s", aurpkg.PackageBase, aurpkg.Version, packageSuffix)

	pkgPath, err = s.controller.CopyFromContainer(containerId, packagePath+pkgName, dest, pkgName)

	if err != nil {
		return "", err
	}

	return pkgPath, nil
}

func (s *AurBuilderService) CleanUpBuildEnv(containerId string) (err error) {
	err = s.controller.RemoveContainerById(containerId)

	if err != nil {
		return err
	}

	return nil
}
