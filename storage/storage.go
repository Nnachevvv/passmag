package storage

import (
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Storage contains user email and - hashed host , encryped password in db
type Storage struct {
	Email       string            `json:"type"`
	Passwords   map[string][]byte `json:"passwords"`
	TimeCreated time.Time         `json:"timecreated"`
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
func StorageFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user : %w", err)
	}

	switch os := runtime.GOOS; os {
	case "windows":
		windowsPath := filepath.Join(usr.HomeDir, "%Appdata%/PasswordManager/vault.bin")
		return filepath.FromSlash(windowsPath), nil
	case "linux":
		return filepath.Join(usr.HomeDir, "/.config", "PasswordManager", "vault.bin"), nil
	case "darwin":
		return filepath.Join("/Library", "Application Support", "PasswordManager", "vault.bin"), nil
	}

	return "", errors.New("current OS is not supported")
}
