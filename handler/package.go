package handler

import (
	"net/http"

	repository "github.com/drzombey/aur-package-builder-api/db/repo"
	"github.com/drzombey/aur-package-builder-api/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var app *model.App

func InitHandlers(a *model.App) {
	app = a

}

func HandleGetPackage(c *gin.Context) {
	repo := repository.PackageRepo{App: app}

	packages, err := repo.GetAurPackages()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occured internally",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, &packages)
}

func HandleAddPackage(c *gin.Context) {
	repo := repository.PackageRepo{App: app}

	var newPackage model.AurPackage

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
