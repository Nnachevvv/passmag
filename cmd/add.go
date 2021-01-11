package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

// the questions to ask
var addQs = []*survey.Question{
	{
		Name:   "host",
		Prompt: &survey.Input{Message: "Enter host adress:"},
	},
	{
		Name:   "password",
		Prompt: &survey.Password{Message: "Enter your password:"},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || len(str) < 8 {
				return errors.New("password should be longer than 8 characters")
			}
			return nil
		},
	},
}

var addCmd = &cobra.Command{

	Use:   "add",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := storage.FilePath()
		if err != nil {
			return err
		}

		if err := storage.VaultExist(path); err != nil {
			return err
		}

		var sessionKey string
		if !viper.IsSet("PASS_SESSION") {
			prompt := &survey.Input{Message: "Please enter your session key :"}
			survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required))
		} else {
			sessionKey = viper.GetString("PASS_SESSION")
		}

		var masterPassword string
		prompt := &survey.Password{Message: "Enter your master password:"}
		survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))

		vaultPwd := argon2.IDKey([]byte(masterPassword), []byte(sessionKey), 1, 64*1024, 4, 32)

		vaultData, err := crypt.DecryptFile(path, vaultPwd)
		if err != nil {
			return err
		}

		err = addPasswords(vaultData, path, vaultPwd)
		if err != nil {
			return err
		}

		fmt.Println("succesfully added")

		return nil
	},
}

func addPasswords(vaultData []byte, path string, vaultPwd []byte) error {
	answers := struct {
		Host     string
		Password string
	}{}

	err := survey.Ask(addQs, &answers)
	if err != nil {
		return err
	}

	s, err := storage.Load(vaultData)
	if err != nil {
		return err
	}

	err = s.Add(answers.Host, answers.Password)
	if err != nil {
		var confirm bool
		editConfirm := &survey.Confirm{Message: "Do you want to edit host with newly password"}
		survey.AskOne(editConfirm, &confirm, survey.WithValidator(survey.Required))
		if confirm {
			s.Edit(answers.Host, answers.Password)
		}
		return nil
	}

	byteData, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal map : %w", err)
	}

	err = crypt.EncryptFile(path, byteData, vaultPwd)

	if err != nil {
		return fmt.Errorf("failed to encrypt sessionData : %w", err)
	}
	return nil
}
