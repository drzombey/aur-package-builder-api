package handler

import (
	"net/http"

	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/builder"
	repository "github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/db/repo"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model"
	"github.com/drzombey/aur-rpc-client-go/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var app *model.App

func InitHandlers(a *model.App) {
	app = a

}

func HandleGetAlreadyBuildPackages(c *gin.Context) {
	repo := repository.PackageRepo{App: app}

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
	repo := repository.PackageRepo{App: app}

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
	repo := repository.PackageRepo{App: app}

	var newPackage types.Package

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
			logrus.Warnf("Warning package already exist Package[id: %s, name: %s, version: %s]", result.ID, result.Name, result.Version)
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "Package already exist, please see the mirror if it's available",
			})
			return
		}
	}

	/*err = repo.AddAurPackage(&newPackage)

	if err != nil {
		logrus.Errorf("Error occured [error: %s]", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}*/

	builder, err := builder.NewAurBuilderService(&app.Config.Docker)

	if err != nil {
		handleError(c, err)
		return
	}

	containerId, err := builder.BuildAurPackage(&newPackage)

	if err != nil {
		handleError(c, err)
		return
	}

	pkgPath, err := builder.CopyPackageToDestination(containerId, "", &newPackage)

	if err != nil {
		handleError(c, err)
		return
	}

	builder.CleanUpBuildEnv(containerId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"msg":       "Package currently creating",
		"processId": containerId,
		"pkgPath":   pkgPath,
	})
}
