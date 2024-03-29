package tasks

import (
	"fmt"

	repository "github.com/drzombey/aur-package-builder-api/cmd/api/db/repo"
	"github.com/drzombey/aur-package-builder-api/pkg/builder"
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
	"github.com/google/uuid"
)

func (at *ApiTask) UpdateAllPackages(taskId uuid.UUID, taskName string) {
	repo, err := repository.NewPackageRepo(at.app.Database)
	if err != nil {
		message := fmt.Sprintf("Cannot establish database connection [Error: %s]", err)
		at.logTaskError(taskId, taskName, message)
		return
	}

	packages, err := repo.GetAlreadyBuildPackages()
	if err != nil {
		message := fmt.Sprintf("Cannot fetch packages [Error: %s]", err)
		at.logTaskError(taskId, taskName, message)
		return
	}

	registryData := docker.RegistryData{
		Image:    at.app.Docker.ContainerImage,
		Username: at.app.Docker.Username,
		Password: at.app.Docker.Password,
		UseAuth:  at.app.Docker.Auth,
	}

	containerController, err := docker.NewContainerController(&registryData)
	if err != nil {
		message := fmt.Sprintf("Cannot initialize container controller [Error: %s]", err)
		at.logTaskError(taskId, taskName, message)
		return
	}

	service, err := builder.NewAurBuilderService(containerController, at.storageProvider)
	if err != nil {
		message := fmt.Sprintf("Cannot initialize service service [Error: %s]", err)
		at.logTaskError(taskId, taskName, message)
		return
	}

	for _, pkg := range *packages {
		aurPkg, err := repo.GetAurPackageInfoByName(pkg.Name)

		if err != nil {
			message := fmt.Sprintf("Cannot get information about package %s by aur rpc call [Error: %s]", pkg.Name, err)
			at.logTaskError(taskId, taskName, message)
			continue
		}

		if aurPkg == nil {
			message := fmt.Sprintf("Package %s with package base id %d not found in aur ", pkg.Name, pkg.PackageBaseID)
			at.logTaskError(taskId, taskName, message)
			continue
		}

		if aurPkg.Version == pkg.Version {
			continue
		}

		err = repo.DeleteAurPackage(pkg)
		if err != nil {
			return
		}

		err = repo.AddAurPackage(*aurPkg)
		if err != nil {
			return
		}

		_, err = service.StartBuildAurPkgRoutine(aurPkg, at.app.PackagePath)
		if err != nil {
			return
		}
	}

	at.logTaskCompleted(taskId, taskName)
}
