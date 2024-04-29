package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	client *mongo.Client
	optns  *options.ClientOptions
}

func NewMongoDBWrapper(ctx context.Context, uri string, poolSize uint64) (*MongoDB, error) {
	newDb := &MongoDB{client: nil, optns: options.Client().ApplyURI(uri).SetMaxPoolSize(poolSize)}
	err := newDb.connect(ctx)
	if err != nil {
		return nil, err
	}

	return newDb, nil
}

func (wrapper *MongoDB) connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, wrapper.optns)

	if err != nil {
		return err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	wrapper.client = client

	return nil
}

func (wrapper *MongoDB) Close(ctx context.Context) error {
	if err := wrapper.client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (wrapper *MongoDB) CreateDocument(ctx context.Context, database, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.InsertOne(context.Background(), document)
	return result, err
}

func (wrapper *MongoDB) CreateDocuments(ctx context.Context, database, collection string, documents []interface{}) (*mongo.InsertManyResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	results, err := coll.InsertMany(context.Background(), documents)
	return results, err
}

func (wrapper *MongoDB) Exists(ctx context.Context, database, collection string, filter interface{}) int64 {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.CountDocuments(context.Background(), filter)

	if err != nil {
		result = 0
	}

	return result
}

func (wrapper *MongoDB) GetDocument(ctx context.Context, database, collection string, filter interface{}) (*mongo.SingleResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result := coll.FindOne(context.Background(), filter)
	return result, result.Err()
}

func (wrapper *MongoDB) GetDocumentWithUpdate(ctx context.Context, database, collection string, filter interface{}, update interface{}) (*mongo.SingleResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result := coll.FindOneAndUpdate(context.Background(), filter, update)
	return result, result.Err()
}

func (wrapper *MongoDB) GetDocuments(ctx context.Context, database, collection string, filter interface{}) (*mongo.Cursor, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.Find(context.Background(), filter)
	return result, err
}

func (wrapper *MongoDB) UpdateDocument(ctx context.Context, database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.UpdateOne(context.Background(), filter, update)
	return result, err
}

func (wrapper *MongoDB) UpdateDocumentByID(ctx context.Context, database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.UpdateByID(context.Background(), filter, update)
	return result, err
}

func (wrapper *MongoDB) UpdateDocuments(ctx context.Context, database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := wrapper.client.Database(database).Collection(collection)
	result, err := coll.UpdateMany(context.Background(), filter, update)
	return result, err
}
