package config

import (
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
	"github.com/drzombey/aur-package-builder-api/pkg/mongo"
)

type AppConfig struct {
	WebserverPort int
	Debug         bool
	LogLevel      string
	JaegerURL     string
	Database      mongo.MongoDbConfig
	Docker        docker.DockerConfig
}
