package storage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Storage contains user email and - hashed host , encryped password in db
type Storage struct {
	Email     string
	Passwords map[string][]byte
}

// New creates a new Storage object from passed email adress and current hashed password in database
func New(data bson.M, email string) Storage {
	s := Storage{Email: email}
	//binaryVaultPwd := data["email"].(primitive.Binary).Data

	for k := range data {
		if k != "_id" && k != "vaultPwd" {
			s.Passwords[string(k)] = data[k].(primitive.Binary).Data
		}
	}
	return s
}
