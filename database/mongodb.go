package database

import (
	"context"
	"fmt"
	"os"

	"github.com/drzombey/aur-package-builder-api/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongoDB(app *model.App) (*mongo.Client, context.Context) {
	log.Info("Setting up mongo database.....")
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		app.Config.Database.User,
		app.Config.Database.Password,
		app.Config.Database.Host,
		app.Config.Database.Port,
	)

	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return client, ctx
}
