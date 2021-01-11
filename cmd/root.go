package cmd

import (
	"context"
	"log"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"

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

func EnterSession() ([]byte, []byte, string, error) {
	path, err := storage.FilePath()
	if err != nil {
		return nil, nil, "", err
	}

	if err := storage.VaultExist(path); err != nil {
		return nil, nil, "", err
	}

	var sessionKey string
	if !viper.IsSet("PASS_SESSION") {
		prompt := &survey.Input{Message: "Please enter your session key :"}
		survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required))
	} else {
		sessionKey = viper.GetString("PASS_SESSION")
	}

	var masterPassword string
	prompt := &survey.Password{Message: "Enter your  master password:"}
	survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))

	vaultPwd := argon2.IDKey([]byte(masterPassword), []byte(sessionKey), 1, 64*1024, 4, 32)

	vaultData, err := crypt.DecryptFile(path, vaultPwd)

	if err != nil {
		return nil, nil, "", err
	}

	return vaultData, vaultPwd, path, err
}
