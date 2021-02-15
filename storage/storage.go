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

	mongoCli "github.com/Nnachevvv/passmag/cmd/mongo"
	"github.com/Nnachevvv/passmag/crypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

//Storage contains user email and - names with password
type Storage struct {
	Email       string            `json:"type"`
	Passwords   map[string][]byte `json:"passwords"`
	TimeCreated time.Time         `json:"timecreated"`
}

// New creates a new Storage object from passed email adress and current encrypted password in database
func New(data map[string]interface{}, creationTime time.Time) (Storage, error) {
	if _, ok := data["email"].(string); !ok {
		return Storage{}, errors.New("failed to get entry only with passwords")
	}

	s := Storage{TimeCreated: creationTime, Email: data["email"].(string)}

	s.Passwords = make(map[string][]byte)

	for k := range data {
		if k != "email" {
			s.Passwords[k] = data[k].(primitive.Binary).Data
		}
	}

	return s, nil
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
		return errors.New("this name already exist")
	}

	s.Passwords[name] = []byte(password)
	return nil
}

// Remove remove name and password from password manager
func (s *Storage) Remove(name string) error {
	if _, ok := s.Passwords[name]; !ok {
		return errors.New("this name not exist in your vault")
	}

	delete(s.Passwords, name)
	return nil
}

// Get gets password if exist from given name
func (s *Storage) Get(name string) (string, error) {
	if _, ok := s.Passwords[name]; !ok {
		return "", errors.New("this name not exist in your vault")
	}

	return string(s.Passwords[name]), nil
}

// ErrCreateUser throw by db when try to insert user
var ErrCreateUser = errors.New("failed to add user to db")

//SyncStorage syncs storage to server if user have connection , otherwise it's throw error
func (s *Storage) SyncStorage(password []byte, mdb mongoCli.Database, client *mongo.Client) error {
	s.TimeCreated = time.Now()
	byteData, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal map : %w", err)
	}

	vaultPwd := argon2.IDKey(password, []byte(s.Email), 1, 64*1024, 4, 32)
	vaultData, err := crypt.Encrypt(byteData, vaultPwd)
	if err != nil {
		return fmt.Errorf("failed to add user to db :%w", err)
	}

	err = mdb.Insert(s.Email, vaultData, client)

	if err != nil {
		return ErrCreateUser
	}

	return nil
}

// Change edits password for given name
func (s *Storage) Change(name, password string) {
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
