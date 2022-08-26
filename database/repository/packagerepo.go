package repository

import (
	"context"

	"github.com/drzombey/aur-package-builder-api/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type PackageRepository struct {
	client   *mongo.Client
	ctx      context.Context
	database string
}

func (r PackageRepository) InitRepo(client *mongo.Client, ctx *context.Context, database string) {
	client = client
	ctx = ctx
	database = database
}

func (r PackageRepository) GetPackageList() []model.AurPackage {

}
