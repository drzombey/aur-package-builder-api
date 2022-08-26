package model

import (
	"github.com/drzombey/aur-package-builder-api/model/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	DbClient mongo.Client
	Config   config.AppConfig
}
