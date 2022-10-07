package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

func (c *ContainerController) EnsureImage() (err error) {

	var authConfig types.AuthConfig
	var authStr string

	if c.RegistryData != nil {
		if c.RegistryData.UseAuth {
			authConfig = types.AuthConfig{
				Username: c.RegistryData.Username,
				Password: c.RegistryData.Password,
			}

			encodedJSON, err := json.Marshal(authConfig)
			if err != nil {
				return err
			}
			authStr = base64.URLEncoding.EncodeToString(encodedJSON)
		}
	}

	logrus.Infof("Pulling new image: %s", c.RegistryData.Image)
	reader, err := c.cli.ImagePull(context.Background(), c.RegistryData.Image, types.ImagePullOptions{RegistryAuth: authStr})

	if err != nil {
		return err
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	return nil
}
