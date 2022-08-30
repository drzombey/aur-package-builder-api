package db

import (
	"context"
	"fmt"

	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model"
	"github.com/drzombey/aur-package-builder-api/cmd/aur_builder_api/model/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct {
	Client *mongo.Client
	app    config.AppConfig
}

func NewMongoStore(app *model.App) (*mongoStore, error) {
	uri := generateUri(app)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	return &mongoStore{Client: client}, nil
}

func generateUri(app *model.App) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		app.Config.Database.User,
		app.Config.Database.Password,
		app.Config.Database.Host,
		app.Config.Database.Port,
	)
}
