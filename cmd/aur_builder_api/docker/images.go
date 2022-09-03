package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

func (c *ContainerController) EnsureImage(image string) (err error) {

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

	reader, err := c.cli.ImagePull(context.Background(), image, types.ImagePullOptions{RegistryAuth: authStr})

	if err != nil {
		return err
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	return nil
}
