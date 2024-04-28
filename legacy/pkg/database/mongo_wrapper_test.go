package database

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	testWrapper *MongoDBWrapper
	testDB      = "testdb"
	testColl    = "testcoll"
)

type Player struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	IdStr      string
	Name       string   `bson:"name,omitempty"`
	Email      string   `bson:"email,omitempty"`
	Tags       []string `bson:"tags,omitempty"`
	Level      int      `bson:"level,omitempty"`
	RandomsTag string
}

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

	docPlayer := &Player{
		Name:       "Tanerius",
		IdStr:      "",
		Email:      "tanerius@player.com",
		Tags:       []string{"a", "b", "c", "d"},
		Level:      1,
		RandomsTag: "Some Tag Does Not get Saved",
	}

	inserted, err := testWrapper.CreateDocument(testDB, testColl, docPlayer)
	if err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	} else {
		if oid, ok := inserted.InsertedID.(primitive.ObjectID); ok {
			docPlayer.IdStr = oid.String()
		} else {
			// Not objectid.ObjectID, do what you want
			t.Fatalf("CreateDocument didnt return an oid")
		}

		if docPlayer.IdStr == "" {
			t.Fatalf("CreateDocument empty string")
		}
	}
}

func TestGetDocument(t *testing.T) {
	clearCollection()

	docPlayer := &Player{
		Name:       "Tanerius",
		IdStr:      "",
		Email:      "tanerius@player.com",
		Tags:       []string{"a", "b", "c", "d"},
		Level:      1,
		RandomsTag: "TestGetDocument",
	}

	testWrapper.CreateDocument(testDB, testColl, docPlayer)

	result, err := testWrapper.GetDocument(testDB, testColl, bson.D{{"email", "tanerius@player.com"}})
	if err != nil {
		t.Fatalf("GetDocument failed: %v", err)
	}

	var retrievedDoc Player
	if err := result.Decode(&retrievedDoc); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if retrievedDoc.Name != "Tanerius" {
		t.Fatalf("Invalid document: %v", retrievedDoc)
	}
}

func TestUpdateDocument(t *testing.T) {
	clearCollection()

	docPlayer := &Player{
		Name:       "Tanerius",
		IdStr:      "",
		Email:      "tanerius@player.com",
		Tags:       []string{"a", "b", "c", "d"},
		Level:      1,
		RandomsTag: "TestUpdateDocument",
	}

	testWrapper.CreateDocument(testDB, testColl, docPlayer)

	update := bson.D{{"$set", bson.D{{"name", "UpdatedTanerius"}}}}
	_, err := testWrapper.UpdateDocument(testDB, testColl, bson.D{{"email", "tanerius@player.com"}}, update)
	if err != nil {
		t.Fatalf("UpdateDocument failed: %v", err)
	}

}
