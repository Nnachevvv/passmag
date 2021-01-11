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

//Storage contains user email and - names with password
type Storage struct {
	Email       string            `json:"type"`
	Passwords   map[string][]byte `json:"passwords"`
	TimeCreated time.Time         `json:"timecreated"`
}

// New creates a new Storage object from passed email adress and current encrypted password in database
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

// Add adds name and password if name doesn't exist
func (s *Storage) Add(name, password string) error {
	if _, ok := s.Passwords[name]; ok {
		return errors.New("this name alredy exist")
	}

	s.Passwords[name] = []byte(password)
	return nil
}

// Remove remove name and password from password manager
func (s *Storage) Remove(name string) error {
	if _, ok := s.Passwords[name]; !ok {
		return errors.New("this name not exist in our db")
	}

	delete(s.Passwords, name)
	return nil
}

// Get gets password if exist from given name
func (s *Storage) Get(name string) (string, error) {
	if _, ok := s.Passwords[name]; !ok {
		return "", errors.New("this name not exist in your password manager ")
	}

	return string(s.Passwords[name]), nil
}

// Edit edits password
func (s *Storage) Edit(name, password string) {
	s.Passwords[name] = []byte(password)
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
