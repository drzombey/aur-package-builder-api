package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

func (c *ContainerController) EnsureImage(image string) (err error) {
	reader, err := c.cli.ImagePull(context.Background(), image, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)
	return nil
}
