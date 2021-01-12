package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Used for flags.
	cfgFile           string
	userLicense       string
	collection        *mongo.Collection
	sessionCollection *mongo.Collection
	ctx               context.Context

	rootCmd = &cobra.Command{
		Use:   "passmag",
		Short: "A password manager used to store securely passwords",
		Long:  `passmag`,
	}
)

func init() {
	rootCmd.AddCommand(login)
	rootCmd.AddCommand(initialize)
	rootCmd.AddCommand(remove)
	rootCmd.AddCommand(add)
	rootCmd.AddCommand(get)
	rootCmd.AddCommand(edit)
	rootCmd.AddCommand(change)
	rootCmd.AddCommand(logout)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("manager")
	collection = db.Collection("users")
	viper.AutomaticEnv()
	//defer client.Disconnect(ctx)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// ErrCreateUser throw by db when try to insert user
var ErrCreateUser = errors.New("failed to add user to db")

// SyncVault syncs current state of vault to password if internet connection is provided
func SyncVault(s storage.Storage, password []byte) error {
	byteData, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal map : %w", err)
	}

	vaultPwd := argon2.IDKey([]byte(s.Email), password, 1, 64*1024, 4, 32)
	byteEncryptedData, err := crypt.Encrypt(byteData, vaultPwd)
	if err != nil {
		return fmt.Errorf("failed to add user to db", err)
	}

	_, err = collection.InsertOne(ctx, bson.D{{Key: "email", Value: s.Email},
		{Key: "vault", Value: byteEncryptedData},
	})

	if err != nil {
		return ErrCreateUser
	}

	return nil
}
