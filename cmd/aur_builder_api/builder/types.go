package builder

import (
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model/config"
)

type AurBuilderService struct {
	cfg config.DockerConfig
}

func NewAurBuilderService(cfg *config.DockerConfig) *AurBuilderService {
	return &AurBuilderService{cfg: *cfg}
}
