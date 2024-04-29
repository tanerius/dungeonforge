package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tanerius/dungeonforge/pkg/config"
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	DisplayName string             `bson:"displayname,omitempty"`
	Email       string             `bson:"email,omitempty"`
	Username    string             `bson:"username,omitempty"`
	Password    string             `bson:"password,omitempty"`
	Created     time.Time          `bson:"created,omitempty"`
	Token       string             `bson:"token,omitempty"`
	Validated   bool               `bson:"validated,omitempty"`
	LastSeen    time.Time          `bson:"lastseen,omitempty"`
	IsOnline    bool               `bson:"online,omitempty"`
	ClientId    string             `bson:"-"`
	log         logging.ILogger    `bson:"-"`
	conf        config.IConfig     `bson:"-"`
	db          *database.MongoDB  `bson:"-"`
}

func (r *User) GetId() string {
	return r.ID.Hex()
}

func (r *User) Write(ctx context.Context) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"displayname": r.DisplayName,
			"email":       r.Email, // TODO: think about not allowing this
			"username":    r.Username,
			"password":    r.Password, // pass should always be hashed
			"token":       r.Token,
			"validated":   r.Validated,
			"lastSeen":    now,
			"online":      r.IsOnline,
		},
	}

	_, err := r.db.UpdateDocumentByID(ctx, "dungeondb", "users", bson.D{{"_id", r.ID}}, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *User) Logout(ctx context.Context) error {
	r.log.LogInfo("entities.User.Logout()")
	if err := r.Write(ctx); err != nil {
		return err
	}

	return nil
}

func (r *User) Validate(ctx context.Context, token string) error {
	r.log.LogInfo("entities.User.Validate()")

	if r.Validated {
		return nil
	}

	if r.Token != token {
		errors.New("invalid token")
	}

	r.Validated = true
	if err := r.Write(ctx); err != nil {
		errors.New("cannot write user")
	}

	return nil
}

// This is basically a login feature for the user
func GetUser(ctx context.Context, db *database.MongoDB, username, password string) (*User, error) {
	logger := logging.NewLogger()
	conf := config.NewIConfig()

	logger.LogInfo("entities.GetUser()")
	// Create the update
	now := time.Now()
	update := bson.M{
		"$set": bson.M{"lastSeen": now, "online": true},
	}
	// TODO: use config to set values
	result, err := db.GetDocumentWithUpdate(ctx, "dungeondb", "users", bson.M{"name": username, "password": database.GetMD5Hash(password)}, update)
	if err != nil {
		return nil, err
	}

	var retrievedDoc *User = &User{}

	if err := result.Decode(retrievedDoc); err != nil {
		return nil, err
	}

	retrievedDoc.LastSeen = now
	retrievedDoc.log = logger
	retrievedDoc.conf = conf
	retrievedDoc.db = db

	return retrievedDoc, nil
}

func RegisterUser(ctx context.Context, db *database.MongoDB, email, username, pass string) (*User, error) {

	logger := logging.NewLogger()
	conf := config.NewIConfig()

	logger.LogInfo("entities.RegisterUser()")

	hashedPass := database.GetMD5Hash(pass)

	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": email},
			bson.M{"username": username},
		},
	}

	count := db.Exists(ctx, "dungeondb", "users", filter)
	if count > 0 {
		err := errors.New("user already exists")
		logger.LogError(err, "user exists")
		return nil, err
	}

	newUser := &User{
		Email:     email,
		Username:  username,
		Password:  hashedPass,
		Created:   time.Now(),
		LastSeen:  time.Now(),
		Token:     uuid.NewString(),
		Validated: false,
		IsOnline:  true,
		conf:      conf,
		db:        db,
		log:       logger,
	}

	inserted, err := db.CreateDocument(ctx, "dungeondb", "users", newUser)
	if err != nil {
		logger.LogError(err, "")
		return nil, err
	} else {
		if oid, ok := inserted.InsertedID.(primitive.ObjectID); ok {
			newUser.ID = oid
		}
		return newUser, nil
	}
}
