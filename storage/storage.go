package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

// Load unmarshal json to Storage struct
func Load(vaultData []byte) (Storage, error) {
	var s Storage

	err := json.Unmarshal(vaultData, &s)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to unmarshal json to struct")
	}

	return s, nil
}

// Add add host and password if host doesn't exist
func (s *Storage) Add(host, password string) error {
	if _, ok := s.Passwords[host]; ok {
		return errors.New("this host alredy exist")
	}

	s.Passwords[host] = []byte(password)
	return nil
}

// Remove remove host and password from password manager
func (s *Storage) Remove(host string) error {
	if _, ok := s.Passwords[host]; ok {
		return errors.New("this host not exist in our db")
	}

	delete(s.Passwords, host)
	return nil
}

// Edit edits host and password
func (s *Storage) Edit(host, password string) error {
	s.Passwords[host] = []byte(password)
	return nil
}

// FilePath returns path for currently used operating system
func FilePath() (string, error) {
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

// VaultExist check if given vault is present in path
func VaultExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.New("please login first")
	}
	return nil
}
