package cmd

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/argon2"
)

// the questions to ask
var loginQs = []*survey.Question{
	{
		Name:   "email",
		Prompt: &survey.Input{Message: "Enter your email adress:"},
		Validate: func(val interface{}) error {
			email, ok := val.(string)
			if !ok || len(email) < 8 {
				return errors.New("email should be longer than 8 characters")
			}

			return nil
		},
	},
	{
		Name:   "masterpassword",
		Prompt: &survey.Password{Message: "Enter your  master password:"},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || len(str) < 8 {
				return errors.New("password should be longer than 8 characters")
			}
			return nil
		},
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to password manager CLI",
	Long:  "login to password manager CLI and seal vault locally with generated random decription key",
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := struct {
			Email          string
			MasterPassword string
		}{}

		err := survey.Ask(loginQs, &answers)
		if err != nil {
			return err
		}

		decryptedVault, err := AuthAndGetVault(answers.Email, answers.MasterPassword)
		if err != nil {
			return err
		}

		sessionKey := make([]byte, 32)
		rand.Read(sessionKey)
		fmt.Println(string("-------------"))
		fmt.Println(string(string(sessionKey)))
		fmt.Println(string(sessionKey))
		fmt.Println(string("-------------"))

		path := storage.OperatingSystem() + "vault.bin"

		vaultPwd := argon2.IDKey([]byte(answers.MasterPassword), sessionKey, 1, 64*1024, 4, 32)
		err = crypt.EncryptFile(path, decryptedVault, vaultPwd)
		if err != nil {
			return fmt.Errorf("failed to encrypt sessionData : %w", err)
		}

		fmt.Println("You're session key is : " + string(sessionKey) + "To unlock your vault\n" +
			"set session key to `PASS_SESSION` enviroment variable like this: \n" +
			"export PASS_SESSION=" + string(sessionKey))

		return nil
	},
}

func AuthAndGetVault(email string, password string) ([]byte, error) {
	db := struct {
		Email string `bson:"email"`
	}{}

	collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&db)
	fmt.Printf("e: %s", db.Email)
	if db.Email == "" || string(email) != db.Email {
		return nil, errors.New("failed to find this account")
	}

	var record bson.M
	collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&record)
	vaultPwd := argon2.IDKey([]byte(password), []byte(email), 1, 64*1024, 4, 32)

	encyrptedVault, err := crypt.Decrypt(record["vault"].(primitive.Binary).Data, vaultPwd)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data value: %w ", err)
	}
	return encyrptedVault, nil

}
