package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
)

var initializeCmd = &cobra.Command{

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

		err = SyncVault(s, []byte(password))
		if err != nil {
			return err
		}

		fmt.Printf("%s is successfully created!\n", email)
		return nil
	},
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
				email, ok := val.(string)
				if !ok || len(email) < 8 {
					return errors.New("email should be longer than 8 characters")
				}

				_, err := service.Find(email)

				if err != nil {
					return fmt.Errorf("email address %s exist in our database", email)
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

	err = survey.Ask(qs, &answers)
	if err != nil {
		return "", "", fmt.Errorf("failed to process input : %w", err)
	}

	if answers.MasterPassword != answers.ConfirmPassword {
		return "", "", errors.New("passwords must match")
	}
	return answers.Email, answers.MasterPassword, nil
}
