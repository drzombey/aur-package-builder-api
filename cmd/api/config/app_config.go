package config

import (
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
	"github.com/drzombey/aur-package-builder-api/pkg/mongo"
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
)

type AppConfig struct {
	WebserverPort   int
	Debug           bool
	LogLevel        string
	JaegerURL       string
	PackagePath     string
	StorageProvider string
	Database        mongo.MongoDbConfig
	Docker          docker.DockerConfig
	Bucket          storage.S3Config
}
