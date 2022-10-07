package repository

import (
	"context"

	"github.com/drzombey/aur-package-builder-api/pkg/aur"
	customMongo "github.com/drzombey/aur-package-builder-api/pkg/mongo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type PackageRepo struct {
	store customMongo.IMongoStore[aur.Package]
}

func NewPackageRepo(config customMongo.MongoDbConfig) (*PackageRepo, error) {
	store, err := customMongo.New[aur.Package](
		config,
		"package",
	)

	if err != nil {
		logrus.Errorf("Failed to initalize mongodb [error: %s]", err)
		return nil, err
	}

	return &PackageRepo{
		store: store,
	}, nil
}

func (pr PackageRepo) GetPackageFromAur(name string) ([]aur.Package, error) {
	response, err := aur.FindPackageByNameInAur(name)
	if err != nil {
		logrus.Errorf("Failed to get packages from arch user repository rpc interface [error: %s]", err)
		return nil, err
	}
	return response.Packages, nil
}

func (pr PackageRepo) GetAlreadyBuildPackageByAurIdAndVersion(id int64, version string) (*aur.Package, error) {
	filter := customMongo.StoreFilter{
		"id":      id,
		"version": version,
	}

	foundPackage, err := pr.store.FindOneBy(context.TODO(), filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logrus.Errorf("Failed to find package from mongodb [error: %s]", err)
		return nil, err
	}

	return foundPackage, nil
}

func (pr PackageRepo) GetAlreadyBuildPackages() (*[]aur.Package, error) {
	result, err := pr.store.Get(context.Background())

	if err != nil {
		logrus.Errorf("Failed to get packages in mongodb [error: %s]", err)
		return nil, err
	}

	return result, nil
}

func (pr PackageRepo) AddAurPackage(model aur.Package) error {
	err := pr.store.Create(context.Background(), model)

	if err != nil {
		logrus.Errorf("Failed to store package in mongodb [error: %s]", err)
		return err
	}

	return nil
}
