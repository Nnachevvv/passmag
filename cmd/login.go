package cmd

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nnachevv/passmag/cmd/mongo"
	"github.com/nnachevv/passmag/random"
	"github.com/nnachevv/passmag/storage"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/argon2"
)

// NewLoinCmd creates a new loginCmd
func NewLoginCmd(md mongo.MongoDatabase) *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "login to password manager CLI",
		Long:  "login to password manager CLI and seal vault locally with generated random description key",
		RunE: func(cmd *cobra.Command, args []string) error {
			email, password, err := loginUserInput()
			if err != nil {
				return err
			}

			decryptedVault, err := getVault(email, password, md)
			if err != nil {
				return err
			}

			sessionKey := random.StringRune(32)

			vaultPwd := argon2.IDKey([]byte(password), []byte(sessionKey), 1, 64*1024, 4, 32)

			path, err := storage.FilePath()
			if err != nil {
				return err
			}

			err = Crypt.EncryptFile(path, decryptedVault, vaultPwd)
			if err != nil {
				return fmt.Errorf("failed to encrypt your vault : %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "You're session key is : "+string(sessionKey)+". To unlock your vault\n"+
				"set session key to `PASS_SESSION` environment variable like this: \n"+
				"export PASS_SESSION="+string(sessionKey))

			return nil
		},
	}
	return loginCmd
}

func getVault(email string, password string, md mongo.MongoDatabase) ([]byte, error) {
	doc, err := md.Find(email)
	if err != nil {
		return nil, err
	}

	vaultPwd := argon2.IDKey([]byte(password), []byte(email), 1, 64*1024, 4, 32)

	encryptedVault, err := Crypt.Decrypt(doc["vault"].(primitive.Binary).Data, vaultPwd)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data value: %w ", err)
	}
	return encryptedVault, nil
}

func loginUserInput() (email string, password string, err error) {
	answers := struct {
		Email          string
		MasterPassword string
	}{}

	loginQs := []*survey.Question{
		{
			Name:   "email",
			Prompt: &survey.Input{Message: "Enter your email address:"},
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

	err = survey.Ask(loginQs, &answers, survey.WithStdio(Stdio.In, Stdio.Out, Stdio.Err))
	if err != nil {
		return "", "", err
	}

	return answers.Email, answers.MasterPassword, nil

}
