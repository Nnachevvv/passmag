package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service contains mongo db client
type Service struct {
	client *mongo.Client
}

// MongoDatabase is used as interface to represent database operation
type MongoDatabase interface {
	Connect()
	Insert(string, []byte) error
	Find(string) (bson.M, error)
}

//Connect connects to mongoDB
func (s *Service) Connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}

	s.client = client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

// Insert wraps InsertOne command from mongodb Driver and insert email and vault to db
func (s *Service) Insert(email string, vault []byte) error {
	db := s.client.Database("manager")
	collection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, bson.D{{Key: "email", Value: email},
		{Key: "vault", Value: vault},
	})

	if err != nil {
		return fmt.Errorf("failed to insert data to db")
	}

	return nil
}

// Find gets data if exist from mongo db client
func (s *Service) Find(email string) (bson.M, error) {
	db := s.client.Database("manager")
	collection := db.Collection("users")

	var record bson.M
	collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&record)

	if record == nil {
		return nil, errors.New("failed to find this account")
	}

	return record, nil
}
