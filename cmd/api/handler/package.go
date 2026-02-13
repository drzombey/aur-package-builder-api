package handler

import (
	"net/http"

	"github.com/drzombey/aur-package-builder-api/cmd/api/config"
	repository "github.com/drzombey/aur-package-builder-api/cmd/api/db/repo"
	"github.com/drzombey/aur-package-builder-api/cmd/api/middleware"
	"github.com/drzombey/aur-package-builder-api/pkg/aur"
	"github.com/drzombey/aur-package-builder-api/pkg/builder"
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler holds all dependencies for the HTTP handlers
type Handler struct {
	config          *config.AppConfig
	storageProvider storage.Provider
	packageRepo     repository.PackageRepository
}

// NewHandler creates a new handler instance with all dependencies
func NewHandler(config *config.AppConfig, storageProvider storage.Provider, packageRepo repository.PackageRepository) *Handler {
	return &Handler{
		config:          config,
		storageProvider: storageProvider,
		packageRepo:     packageRepo,
	}
}

func (h *Handler) HandleGetAlreadyBuildPackages(c *gin.Context) {
	packages, err := h.packageRepo.GetAlreadyBuildPackages()

	if err != nil {
		middleware.AbortWithError(c, middleware.NewInternalServerError("Failed to get packages", err))
		return
	}

	c.IndentedJSON(http.StatusOK, &packages)
}

func (h *Handler) HandleGetAurPackageByName(c *gin.Context) {
	packagename, isPackageNameSet := c.GetQuery("packageName")

	if !isPackageNameSet {
		middleware.AbortWithError(c, middleware.NewBadRequestError("Query parameter packageName not set", nil))
		return
	}

	response, err := h.packageRepo.GetPackageFromAur(packagename)

	if err != nil {
		middleware.AbortWithError(c, middleware.NewInternalServerError("Failed to get packages from AUR", err))
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (h *Handler) HandleBuildPackage(c *gin.Context) {
	var newPackage aur.Package

	if err := c.BindJSON(&newPackage); err != nil {
		middleware.AbortWithError(c, middleware.NewBadRequestError("Invalid JSON structure", err))
		return
	}

	result, err := h.packageRepo.GetAlreadyBuildPackageByAurIdAndVersion(newPackage.ID, newPackage.Version)
	if err != nil {
		middleware.AbortWithError(c, middleware.NewInternalServerError("Failed to check existing package", err))
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

	err = h.packageRepo.AddAurPackage(newPackage)
	if err != nil {
		middleware.AbortWithError(c, middleware.NewInternalServerError("Failed to add package to database", err))
		return
	}

	registryData := docker.RegistryData{
		Image:    h.config.Docker.ContainerImage,
		Username: h.config.Docker.Username,
		Password: h.config.Docker.Password,
		UseAuth:  h.config.Docker.Auth,
	}

	containerController, err := docker.NewContainerController(&registryData)
	if err != nil {
		middleware.AbortWithError(c, middleware.NewInternalServerError("Failed to create container controller", err))
		return
	}

	builderService, err := builder.NewAurBuilderService(containerController, h.storageProvider)
	if err != nil {
		middleware.AbortWithError(c, middleware.NewInternalServerError("Failed to create builder service", err))
		return
	}

	containerId, err := builderService.StartBuildAurPkgRoutine(&newPackage, h.config.PackagePath)

	if err != nil {
		middleware.AbortWithError(c, middleware.NewInternalServerError("Failed to start build routine", err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"msg":       "Package currently creating",
		"processId": containerId,
	})
}
