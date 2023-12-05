package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBWrapper struct {
	client *mongo.Client
}

func NewMongoDBWrapper(uri string, poolSize uint64) *MongoDBWrapper {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(poolSize)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	return &MongoDBWrapper{client: client}
}

func (wrapper *MongoDBWrapper) Close() {
	if err := wrapper.client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (wrapper *MongoDBWrapper) CreateDocument(database, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.InsertOne(context.Background(), document)
	return result, err
}

func (wrapper *MongoDBWrapper) CreateDocuments(database, collection string, documents []interface{}) (*mongo.InsertManyResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	results, err := coll.InsertMany(context.Background(), documents)
	return results, err
}

func (wrapper *MongoDBWrapper) GetDocument(database, collection string, filter interface{}) (*mongo.SingleResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result := coll.FindOne(context.Background(), filter)
	return result, result.Err()
}

func (wrapper *MongoDBWrapper) GetDocuments(database, collection string, filter interface{}) (*mongo.Cursor, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.Find(context.Background(), filter)
	return result, err
}

func (wrapper *MongoDBWrapper) UpdateDocument(database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.UpdateOne(context.Background(), filter, update)
	return result, err
}

func (wrapper *MongoDBWrapper) UpdateDocuments(database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.UpdateMany(context.Background(), filter, update)
	return result, err
}
