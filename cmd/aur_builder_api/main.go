package main

import (
	"fmt"

	docker "github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/container"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/handler"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	app model.App
)

func main() {
	controller, err := docker.NewContainerController()
	if err != nil {
		panic(err)
	}

	//container, _ := controller.ContainerById("457e5bcc644dcb0ff2612c10b4d1a55b01f20e4615f046675894493fee56bb3f")

	//fmt.Print(container.Names)
	//fmt.Print(container.State)

	containers, _ := controller.ContainersByImage("mongo:latest")

	fmt.Println(containers)

	controller.CleanContainersByImage("mongo:latest")

	containers, _ = controller.ContainersByImage("mongo:latest")

	fmt.Println(containers)

	//setupLogFormatter()
	//loadConfig()
	//setupWebserver()
}

func setupWebserver() {
	gin.SetMode(gin.ReleaseMode)

	if app.Config.Debug {
		gin.SetMode(gin.DebugMode)
	}

	server := gin.Default()
	registerHandlers(server)
	server.Run(fmt.Sprintf(":%d", app.Config.WebserverPort))
}

func setupLogFormatter() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func setupDatabase() {
	//store, err := db.NewMongoStore(&app)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("webserverPort", 8080)
	viper.SetDefault("webserverMode", "production")
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

	err := viper.Unmarshal(&app.Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logLevel, err := log.ParseLevel(app.Config.LogLevel)
	if err == nil {
		log.SetLevel(logLevel)
	}

	log.Info("Config loaded.")
}

func registerHandlers(s *gin.Engine) {
	handler.InitHandlers(&app)

	version1 := "/api/v1"

	s.GET(fmt.Sprintf("%s/package", version1), handler.HandleGetPackageList)
	s.POST(fmt.Sprintf("%s/package", version1), handler.HandleAddPackage)
}
