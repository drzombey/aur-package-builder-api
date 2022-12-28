package tasks

import (
	"github.com/drzombey/aur-package-builder-api/cmd/api/config"
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
)

type ApiTask struct {
	app             config.AppConfig
	storageProvider storage.Provider
}

func NewApiTask(app config.AppConfig, provider storage.Provider) *ApiTask {
	return &ApiTask{
		app:             app,
		storageProvider: provider,
	}
}
