package docker

type DockerConfig struct {
	Auth           bool
	Username       string
	Password       string
	ContainerImage string
}
