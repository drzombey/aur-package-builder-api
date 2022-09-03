package config

type AppConfig struct {
	WebserverPort int
	Debug         bool
	LogLevel      string
	Database      MongoDbConfig
	Docker        DockerConfig
}
