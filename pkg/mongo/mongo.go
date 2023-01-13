package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore[T any] struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type StoreFilter = bson.M

type IMongoStore[T any] interface {
	FindOneBy(ctx context.Context, filter interface{}) (*T, error)
	FindBy(ctx context.Context, filter interface{}) (*[]T, error)
	Get(ctx context.Context) (*[]T, error)
	Delete(ctx context.Context, filter interface{}) error
	Create(ctx context.Context, obj T) error
}

func New[T any](config MongoDbConfig, collection string) (IMongoStore[T], error) {
	uri := generateUri(&config)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	cl := client.Database(config.Name).Collection(collection)

	collections, err := cl.Database().ListCollectionNames(context.Background(), bson.D{})

	if err != nil {
		return nil, err
	}

	var found = false
	for _, name := range collections {
		if name == collection {
			found = true
		}
	}

	if !found {
		err := cl.Database().CreateCollection(context.Background(), collection)

		if err != nil {
			return nil, err
		}
	}

	return MongoStore[T]{
		client:     client,
		collection: cl,
	}, nil
}

func generateUri(config *MongoDbConfig) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
}

func (s MongoStore[T]) FindOneBy(ctx context.Context, filter interface{}) (*T, error) {
	var dbo T
	err := s.collection.FindOne(context.TODO(), filter).Decode(&dbo)

	if err != nil {
		return nil, err
	}

	return &dbo, nil
}

func (s MongoStore[T]) FindBy(ctx context.Context, filter interface{}) (*[]T, error) {
	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var result []T
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var entry T
		if err = cur.Decode(&entry); err != nil {
			return nil, err
		}

		result = append(result, entry)
	}

	return &result, nil
}

func (s MongoStore[T]) Get(ctx context.Context) (*[]T, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var result []T
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var entry T
		if err = cur.Decode(&entry); err != nil {
			return nil, err
		}

		result = append(result, entry)
	}

	return &result, nil
}

func (s MongoStore[T]) Create(ctx context.Context, obj T) error {
	_, err := s.collection.InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

func (s MongoStore[T]) Delete(ctx context.Context, filter interface{}) error {
	_, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
