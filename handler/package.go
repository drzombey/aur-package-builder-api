package handler

import (
	"github.com/drzombey/aur-package-builder-api/model"
	"github.com/gin-gonic/gin"
)

var packages = []model.AurPackage{
	{ID: "1", Name: "Test", Version: "Me", CreationDate: "5"},
	{ID: "2", Name: "Test", Version: "Me", CreationDate: "5"},
	{ID: "3", Name: "Test", Version: "Me", CreationDate: "5"},
	{ID: "4", Name: "Test", Version: "Me", CreationDate: "5"},
}

func HandleGetPackage(c *gin.Context) {

}
