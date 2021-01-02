package storage

import (
	"go.mongodb.org/mongo-driver/bson"
)

//Storage contains user email and - hashed host , encryped password in db
type Storage struct {
	Email     string
	Passwords map[string]string
}

// New creates a new Storage object from passed email adress and current hashed password in database
func New(data bson.M, email string) Storage {
	s := Storage{Email: email}
	//binaryVaultPwd := data["email"].(primitive.Binary).Data
	s.Passwords = make(map[string]string)
	for k := range data {
		if k != "_id" && k != "email" {
			s.Passwords[k] = data[k].(string)
		}
	}
	return s
}
