package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

// the questions to ask
var editPwd = []*survey.Question{
	{
		Name:   "host",
		Prompt: &survey.Input{Message: "Enter host for which you want to edit your password adress:"},
	},
	{
		Name:   "password",
		Prompt: &survey.Password{Message: "Enter new password:"},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || len(str) < 8 {
				return errors.New("password should be longer than 8 characters")
			}
			return nil
		},
	},
}

var editCmd = &cobra.Command{

	Use:   "edit",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//TODO check if vault is present?test this
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
		path, err := storage.FilePath()
		if err != nil {
			return err
		}

		vaultData, err := crypt.DecryptFile(path, vaultPwd)

		if err != nil {
			return err
		}

		var s storage.Storage

		err = json.Unmarshal(vaultData, &s)
		if err != nil {
			return err
		}

		answers := struct {
			Host     string
			Password string
		}{}

		err = survey.Ask(editPwd, &answers)
		if err != nil {
			return err
		}

		if _, ok := s.Passwords[answers.Host]; !ok {
			return errors.New("failed to find this host")
		}

		s.Edit(answers.Host, answers.Password)
		byteData, err := json.Marshal(s)
		if err != nil {
			return fmt.Errorf("failed to marshal map : %w", err)
		}

		err = crypt.EncryptFile(path, byteData, vaultPwd)

		if err != nil {
			return fmt.Errorf("failed to encrypt sessionData : %w", err)
		}
		return nil
	},
}
