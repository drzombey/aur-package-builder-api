package repository

import (
	"context"

	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/aur"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/db"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model"
	"github.com/drzombey/aur-rpc-client-go/types"
	"go.mongodb.org/mongo-driver/bson"
)

type PackageRepo struct {
	App *model.App
}

func (pr PackageRepo) GetPackageFromAur(name string) ([]types.Package, error) {
	response, err := aur.FindPackageByNameInAur(name)
	if err != nil {
		return nil, err
	}
	return response.Packages, nil
}

func (pr PackageRepo) GetAlreadyBuildPackages() ([]*model.AurPackage, error) {
	store, err := db.NewMongoStore(pr.App)
	packages := []*model.AurPackage{}

	if err != nil {
		return nil, err
	}

	cur, err := store.Client.Database(pr.App.Config.Database.Name).Collection("packages").Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, err
	}

	if err := cur.All(context.TODO(), &packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func (pr PackageRepo) AddAurPackage(model *model.AurPackage) error {
	store, err := db.NewMongoStore(pr.App)
	if err != nil {
		return err
	}

	col := store.Client.Database(pr.App.Config.Database.Name).Collection("packages")

	_, err = col.InsertOne(context.TODO(), model)

	if err != nil {
		return err
	}

	return nil
}
