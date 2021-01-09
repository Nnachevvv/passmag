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
var getQs = []*survey.Question{
	{
		Name:   "host",
		Prompt: &survey.Input{Message: "Enter host for which you want to get password:"},
	},
}

var getCmd = &cobra.Command{

	Use:   "get",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//TODO check if vault is present?
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
		path, err := storage.StorageFilePath()
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
			Host string
		}{}

		err = survey.Ask(getQs, &answers)
		if err != nil {
			return err
		}

		if _, ok := s.Passwords[answers.Host]; !ok {
			return errors.New("failed to find this password")
		}

		fmt.Println(string(s.Passwords[answers.Host]))

		return nil
	},
}
