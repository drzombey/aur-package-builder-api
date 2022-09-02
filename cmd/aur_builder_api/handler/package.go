package handler

import (
	"net/http"

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

func HandleGetPackageList(c *gin.Context) {
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

func HandleAddPackage(c *gin.Context) {
	repo := repository.PackageRepo{App: app}

	var newPackage types.Package

	if err := c.BindJSON(&newPackage); err != nil {
		logrus.Errorf("Error occured [error: %s]", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "Invalid JSON structure",
		})
		return
	}

	err := repo.AddAurPackage(&newPackage)

	if err != nil {
		logrus.Errorf("Error occured [error: %s]", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"msg":    "Package created",
	})
}
