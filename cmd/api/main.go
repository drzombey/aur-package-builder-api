package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/drzombey/aur-package-builder-api/cmd/api/config"
	repository "github.com/drzombey/aur-package-builder-api/cmd/api/db/repo"
	"github.com/drzombey/aur-package-builder-api/cmd/api/handler"
	"github.com/drzombey/aur-package-builder-api/cmd/api/middleware"
	"github.com/drzombey/aur-package-builder-api/cmd/api/tasks"
	"github.com/drzombey/aur-package-builder-api/pkg/scheduler"
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
	"github.com/drzombey/aur-package-builder-api/pkg/tracing"
)

var (
	app           config.AppConfig
	configPath    string
	taskScheduler *scheduler.TasksScheduler
	packageRepo   *repository.PackageRepo
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

	// Initialize package repository once
	var err error
	packageRepo, err = repository.NewPackageRepo(app.Database)
	if err != nil {
		log.Fatalf("Failed to initialize package repository: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	if app.Debug {
		gin.SetMode(gin.DebugMode)
	}

	server := gin.Default()
	server.Use(middleware.ErrorHandler())
	server.Use(otelgin.Middleware("aur-package-builder-api"))
	registerHandlers(server)
	initBackgroundTasks()

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.WebserverPort),
		Handler: server,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Starting server on port %d", app.WebserverPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Give the server 30 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown error: %v", err)
	}

	// Cleanup resources
	cleanup()

	log.Info("Server shutdown complete")
}

func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)

	// Enable environment variables with prefix
	viper.SetEnvPrefix("AUR_BUILDER")
	viper.AutomaticEnv()

	// Replace dots and dashes in env variables with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.SetDefault("webserverPort", 8080)
	viper.SetDefault("debug", false)
	viper.SetDefault("logLevel", "info")
	viper.SetDefault("jaegerURL", "http://localhost:14268/api/traces")
	viper.SetDefault("packagePath", ".")
	viper.SetDefault("storageProvider", "")
	
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 27017)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "example")
	viper.SetDefault("database.name", "aur_package_builder")

	// Docker defaults  
	viper.SetDefault("docker.containerImage", "")
	viper.SetDefault("docker.username", "")
	viper.SetDefault("docker.password", "")
	viper.SetDefault("docker.auth", false)

	// S3 bucket defaults
	viper.SetDefault("bucket.name", "")
	viper.SetDefault("bucket.region", "")
	viper.SetDefault("bucket.accessKey", "")
	viper.SetDefault("bucket.secretKey", "")

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Error reading config file, using default values and environment variables. %s", err)
	}

	err := viper.Unmarshal(&app)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logLevel, err := log.ParseLevel(app.LogLevel)
	if err == nil {
		log.SetLevel(logLevel)
	}

	log.Info("Config loaded from file and environment variables.")
}

func registerHandlers(s *gin.Engine) {
	// Create handler with dependencies
	handlerInstance := handler.NewHandler(&app, getStorageProvider(), packageRepo)

	version1 := "/api/v1"

	s.GET(fmt.Sprintf("%s/build/package", version1), handlerInstance.HandleGetAlreadyBuildPackages)
	s.POST(fmt.Sprintf("%s/build/package", version1), handlerInstance.HandleBuildPackage)
	s.GET(fmt.Sprintf("%s/aurpackage", version1), handlerInstance.HandleGetAurPackageByName)
}

func cleanup() {
	log.Info("Starting cleanup...")

	// Stop task scheduler if running
	if taskScheduler != nil {
		log.Info("Stopping task scheduler...")
		taskScheduler.ShutdownAllTasks()
	}

	// Add any additional cleanup here (e.g., stopping Docker containers)
	log.Info("Cleanup completed")
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
