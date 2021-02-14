package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/argon2"
)

// NewInitCmd creates a new i
func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize email, password and master password for your password manager",
		Long:  `Set master password`,
		RunE: func(cmd *cobra.Command, args []string) error {
			email, password, err := initUserInput()
			if err != nil {
				return err
			}

			s, err := storage.New(bson.M{"email": email}, time.Now())
			if err != nil {
				return err
			}

			byteData, err := json.Marshal(s)
			if err != nil {
				return fmt.Errorf("failed to marshal map : %w", err)
			}

			vaultPwd := argon2.IDKey([]byte(password), []byte(s.Email), 1, 64*1024, 4, 32)
			vaultData, err := crypt.Encrypt(byteData, vaultPwd)
			if err != nil {
				return fmt.Errorf("failed to add user to db :%w", err)
			}

			err = MongoDB.Insert(s.Email, vaultData, Client)

			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s is successfully created!\n", email)
			return nil
		},
	}

	return initCmd
}

func initUserInput() (email string, password string, err error) {
	answers := struct {
		Email           string
		MasterPassword  string
		ConfirmPassword string
	}{}

	qs := []*survey.Question{
		{
			Name:   "email",
			Prompt: &survey.Input{Message: "Enter your email address:"},
			Validate: func(val interface{}) error {

				if input, ok := val.(string); !ok || len(input) < 8 {
					return errors.New("email should be longer than 8 characters")
				}
				if _, err := MongoDB.Find(val.(string), Client); err == nil {
					return errors.New("email address already exist in our database")
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
	err = survey.Ask(qs, &answers, survey.WithStdio(Stdio.In, Stdio.Out, Stdio.Err))
	if err != nil {
		return "", "", fmt.Errorf("failed to process input : %w", err)
	}

	if answers.MasterPassword != answers.ConfirmPassword {
		return "", "", errors.New("passwords must match")
	}

	return answers.Email, answers.MasterPassword, nil
}
