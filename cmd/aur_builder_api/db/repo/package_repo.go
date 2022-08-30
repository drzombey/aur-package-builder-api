package repository

import (
	"context"

	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/db"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type PackageRepo struct {
	App *model.App
}

func (pr PackageRepo) GetAurPackages() ([]*model.AurPackage, error) {
	store, err := db.NewMongoStore(pr.App)
	packages := []*model.AurPackage{}

	if err != nil {
		logrus.Errorf("Error during mongo client creation [error: %s]", err)
		return nil, err
	}

	cur, err := store.Client.Database(pr.App.Config.Database.Name).Collection("packages").Find(context.TODO(), bson.D{})

	if err != nil {
		logrus.Errorf("Error during find of packages [error: %s]", err)
		return nil, err
	}

	if err := cur.All(context.TODO(), &packages); err != nil {
		logrus.Errorf("Error during cursor search of packages [error: %s]", err)
		return nil, err
	}

	return packages, nil
}

func (pr PackageRepo) AddAurPackage(model *model.AurPackage) error {
	store, err := db.NewMongoStore(pr.App)
	if err != nil {
		logrus.Errorf("Error during mongo client creation [error: %s]", err)
		return err
	}

	col := store.Client.Database(pr.App.Config.Database.Name).Collection("packages")

	_, err = col.InsertOne(context.TODO(), model)

	if err != nil {
		logrus.Errorf("Error occured during insertion of package [error: %s]", err)
		return err
	}

	return nil
}
