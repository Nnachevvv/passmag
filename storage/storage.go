package storage

import (
	"path/filepath"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Storage contains user email and - hashed host , encryped password in db
type Storage struct {
	Email       string
	Passwords   map[string][]byte
	TimeCreated time.Time
}

// New creates a new Storage object from passed email adress and current hashed password in database
func New(data map[string]interface{}) Storage {
	s := Storage{TimeCreated: time.Now(), Email: data["email"].(string)}

	s.Passwords = make(map[string][]byte)
	for k := range data {
		if k != "email" {
			s.Passwords[k] = data[k].(primitive.Binary).Data
		}
	}
	return s
}

//OperatingSystem returns path for currently used operating system
func OperatingSystem() string {
	switch os := runtime.GOOS; os {
	case "windows":
		return filepath.FromSlash("%Appdata%/PasswordManager")
	case "linux":
		return filepath.Join("/.config", "PasswordManager")
	case "darwin":
		return filepath.Join("/Library", "Application Support", "PasswordManager")
	}
	return ""
}
