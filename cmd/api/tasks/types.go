package tasks

import "github.com/drzombey/aur-package-builder-api/cmd/api/config"

type ApiTask struct {
	app config.AppConfig
}

func NewApiTask(app config.AppConfig) *ApiTask {
	return &ApiTask{
		app: app,
	}
}
