package core

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ckalagara/group-a-accounts/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "group-a"
	collectionName = "accounts"
)

type Store interface {
	Update(a model.Account) (model.Account, error)
	Get(id string) (model.Account, error)
	Delete(id string) error
}

type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoStore(ctx context.Context, mongoURI string) *MongoStore {
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to create Mongo client: %v", err)
		return nil
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil
	}

	itemsCollection := client.Database(databaseName).Collection(collectionName)

	log.Println("Successfully connected to MongoDB")

	return &MongoStore{
		client:     client,
		collection: itemsCollection,
	}
}

func (m *MongoStore) Update(a model.Account) (model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.Updated = time.Now()
	filter := bson.M{"_id": a.ID}
	update := bson.M{"$set": a}

	opts := options.Update().SetUpsert(true)
	_, err := m.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return model.Account{}, err
	}
	return a, nil
}

func (m *MongoStore) Get(id string) (model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result model.Account
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Account{}, errors.New("account not found")
		}
		return model.Account{}, err
	}
	return result, nil
}

func (m *MongoStore) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("account not found")
	}
	return nil
}
