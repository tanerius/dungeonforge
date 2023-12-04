package database

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	testWrapper *MongoDBWrapper
	testDB      = "testdb"
	testColl    = "testcoll"
)

func TestMain(m *testing.M) {
	testWrapper = NewMongoDBWrapper("mongodb://dungeonmaster:m123123123@localhost:27017/", 10)
	defer testWrapper.Close()

	// Drop test database after tests
	defer func() {
		if err := testWrapper.client.Database(testDB).Drop(context.Background()); err != nil {
			panic(err)
		}
	}()

	os.Exit(m.Run())
}

// Helper function to clear collection
func clearCollection() {
	if err := testWrapper.client.Database(testDB).Collection(testColl).Drop(context.Background()); err != nil {
		panic(err)
	}
}

func TestCreateDocument(t *testing.T) {
	clearCollection()

	doc := bson.D{{"name", "John Doe"}, {"age", 30}}
	_, err := testWrapper.CreateDocument(testDB, testColl, doc)
	if err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	}

	// TODO: verify the document was correctly inserted
}

func TestGetDocument(t *testing.T) {
	clearCollection()

	doc := bson.D{{"name", "Jane Doe"}, {"age", 25}}
	testWrapper.CreateDocument(testDB, testColl, doc)

	result, err := testWrapper.GetDocument(testDB, testColl, bson.D{{"name", "Jane Doe"}})
	if err != nil {
		t.Fatalf("GetDocument failed: %v", err)
	}

	var retrievedDoc bson.D
	if err := result.Decode(&retrievedDoc); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// TODO: verify the retrieved document is correct
}

func TestUpdateDocument(t *testing.T) {
	clearCollection()

	doc := bson.D{{"name", "John Doe"}, {"age", 30}}
	testWrapper.CreateDocument(testDB, testColl, doc)

	update := bson.D{{"$set", bson.D{{"age", 31}}}}
	_, err := testWrapper.UpdateDocument(testDB, testColl, bson.D{{"name", "John Doe"}}, update)
	if err != nil {
		t.Fatalf("UpdateDocument failed: %v", err)
	}

	// TODO: verify the document was correctly updated
}
