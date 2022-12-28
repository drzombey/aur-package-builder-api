package main

import (
	"flag"
	"fmt"

	"github.com/drzombey/aur-package-builder-api/cmd/api/config"
	"github.com/drzombey/aur-package-builder-api/cmd/api/handler"
	"github.com/drzombey/aur-package-builder-api/pkg/tracing"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	app        config.AppConfig
	configPath string
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
	handler.InitHandlers(&app)

	version1 := "/api/v1"

	s.GET(fmt.Sprintf("%s/build/package", version1), handler.HandleGetAlreadyBuildPackages)
	s.POST(fmt.Sprintf("%s/build/package", version1), handler.HandleBuildPackage)
	s.GET(fmt.Sprintf("%s/aurpackage", version1), handler.HandleGetAurPackageByName)
}
