package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/argon2"
)

var qs = []*survey.Question{
	{
		Name:   "email",
		Prompt: &survey.Input{Message: "Enter your email adress:"},
		Validate: func(val interface{}) error {
			email, ok := val.(string)
			if !ok || len(email) < 8 {
				return errors.New("email should be longer than 8 characters")
			}
			var record bson.M

			collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&record)
			if record != nil {
				return fmt.Errorf("email adress %s exist in our database", email)
			}
			return nil
		},
	},
	{
		Name:   "masterpassword",
		Prompt: &survey.Password{Message: "Enter your password:"},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || len(str) < 8 {
				return errors.New("password should be longer than 8 characters")
			}
			return nil
		},
	},
	{
		Name:   "confirmpassword",
		Prompt: &survey.Password{Message: "Enter again your password:"},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || len(str) < 8 {
				return errors.New("password should be longer than 8 characters")
			}
			return nil
		},
	},
}

var initialize = &cobra.Command{

	Use:   "init",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := struct {
			Email           string
			MasterPassword  string
			ConfirmPassword string
		}{}

		survey.Ask(qs, &answers)

		if answers.MasterPassword != answers.ConfirmPassword {
			return errors.New("passwords doesn't match")
		}

		s := storage.New(bson.M{"email": answers.Email})

		byteData, err := json.Marshal(s)
		if err != nil {
			return fmt.Errorf("failed to marshal map : %w", err)
		}

		vaultPwd := argon2.IDKey([]byte(answers.MasterPassword), []byte(answers.Email), 1, 64*1024, 4, 32)
		byteEncryptedData, err := crypt.Encrypt(byteData, vaultPwd)
		if err != nil {
			return fmt.Errorf("failed to encrypt your data: %w", err)
		}

		_, err = collection.InsertOne(ctx, bson.D{{Key: "email", Value: answers.Email},
			{Key: "vault", Value: byteEncryptedData},
		})

		if err != nil {
			return fmt.Errorf("failed to create user : %w ", err)
		}

		fmt.Printf("%s is successfully created!\n", answers.Email)
		return nil
	},
}
