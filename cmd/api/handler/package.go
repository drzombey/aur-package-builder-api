package handler

import (
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
	"net/http"

	"github.com/drzombey/aur-package-builder-api/cmd/api/config"
	repository "github.com/drzombey/aur-package-builder-api/cmd/api/db/repo"
	"github.com/drzombey/aur-package-builder-api/pkg/aur"
	"github.com/drzombey/aur-package-builder-api/pkg/builder"
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	app             *config.AppConfig
	storageProvider storage.Provider
)

func InitHandlers(a *config.AppConfig, sp storage.Provider) {
	app = a
	storageProvider = sp
}

func HandleGetAlreadyBuildPackages(c *gin.Context) {
	repo, err := repository.NewPackageRepo(app.Database)

	if err != nil {
		logrus.Errorf("Failed to initialize database connection [error: %s]", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}

	packages, err := repo.GetAlreadyBuildPackages()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, &packages)
}

func HandleGetAurPackageByName(c *gin.Context) {
	repo, err := repository.NewPackageRepo(app.Database)

	if err != nil {
		logrus.Errorf("Failed to initialize database connection [error: %s]", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}

	packagename, isPackageNameSet := c.GetQuery("packageName")

	if !isPackageNameSet {
		logrus.Error("Error query parameter not set")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "Query parameter packageName not set!",
		})
		return
	}

	response, err := repo.GetPackageFromAur(packagename)

	if err != nil {
		logrus.Errorf("Error during getting of aur packages [error %s]", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func HandleBuildPackage(c *gin.Context) {
	repo, err := repository.NewPackageRepo(app.Database)

	if err != nil {
		handleError(c, err)
		return
	}

	var newPackage aur.Package

	if err := c.BindJSON(&newPackage); err != nil {
		handleInvalidJsonStructure(c, err)
		return
	}

	result, err := repo.GetAlreadyBuildPackageByAurIdAndVersion(newPackage.ID, newPackage.Version)
	if err != nil {
		handleError(c, err)
		return
	}

	if result != nil {
		if result.Version == newPackage.Version {
			logrus.Warnf("Warning package already exist Package[id: %d, name: %s, version: %s]", result.ID, result.Name, result.Version)
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "Package already exist, please see the mirror if it's available",
			})
			return
		}
	}

	err = repo.AddAurPackage(newPackage)
	if err != nil {
		logrus.Errorf("Error occured [error: %s]", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}

	registryData := docker.RegistryData{
		Image:    app.Docker.ContainerImage,
		Username: app.Docker.Username,
		Password: app.Docker.Password,
		UseAuth:  app.Docker.Auth,
	}

	containerController, err := docker.NewContainerController(&registryData)
	if err != nil {
		handleError(c, err)
		return
	}

	builderService, err := builder.NewAurBuilderService(containerController, storageProvider)
	if err != nil {
		handleError(c, err)
		return
	}

	containerId, err := builderService.StartBuildAurPkgRoutine(&newPackage, app.PackagePath)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"msg":       "Package currently creating",
		"processId": containerId,
	})
}
