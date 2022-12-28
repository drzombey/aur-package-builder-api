package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/drzombey/aur-package-builder-api/cmd/api/config"
	"github.com/drzombey/aur-package-builder-api/cmd/api/handler"
	"github.com/drzombey/aur-package-builder-api/cmd/api/tasks"
	"github.com/drzombey/aur-package-builder-api/pkg/scheduler"
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
	"github.com/drzombey/aur-package-builder-api/pkg/tracing"
)

var (
	app           config.AppConfig
	configPath    string
	taskScheduler *scheduler.TasksScheduler
)

func init() {
	flag.StringVar(&configPath, "config_path", ".", "path to search for a config.yaml")
	flag.Parse()
}

func main() {
	setupLogging()
	loadConfig()

	closer := tracing.Setup(app.JaegerURL)
	defer closer()

	gin.SetMode(gin.ReleaseMode)

	if app.Debug {
		gin.SetMode(gin.DebugMode)
	}

	server := gin.Default()
	server.Use(otelgin.Middleware("aur-package-builder-api"))
	registerHandlers(server)
	initBackgroundTasks()
	server.Run(fmt.Sprintf(":%d", app.WebserverPort))
}

func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)

	viper.SetDefault("webserverPort", 8080)
	viper.SetDefault("webserverMode", "production")
	viper.SetDefault("jaegerURL", "http://localhost:14268/api/traces")
	viper.SetDefault("packagePath", ".")
	viper.SetDefault("storageProvider", "")
	viper.SetDefault("database", map[string]interface{}{
		"host":     "localhost",
		"port":     27017,
		"user":     "root",
		"password": "example",
		"name":     "aur_package_builder",
	})

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Error reading config file, using default values. %s", err)
	}

	err := viper.Unmarshal(&app)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logLevel, err := log.ParseLevel(app.LogLevel)
	if err == nil {
		log.SetLevel(logLevel)
	}

	log.Info("Config loaded.")
}

func registerHandlers(s *gin.Engine) {
	handler.InitHandlers(&app, getStorageProvider())

	version1 := "/api/v1"

	s.GET(fmt.Sprintf("%s/build/package", version1), handler.HandleGetAlreadyBuildPackages)
	s.POST(fmt.Sprintf("%s/build/package", version1), handler.HandleBuildPackage)
	s.GET(fmt.Sprintf("%s/aurpackage", version1), handler.HandleGetAurPackageByName)
}

func getStorageProvider() storage.Provider {
	switch app.StorageProvider {
	case "object":
		return storage.NewS3Provider(&app.Bucket)
	default:
		return storage.NewFilesystemProvider(app.PackagePath)
	}
}

func initBackgroundTasks() {
	taskScheduler = scheduler.NewTasksScheduler()
	apiTask := tasks.NewApiTask(app, getStorageProvider())
	taskScheduler.ScheduleTask(apiTask.UpdateAllPackages, 3600, "UpdatePackagesTask")
}
