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

var get = &cobra.Command{

	Use:   "get",
	Short: "Get password from your vault",
	Long:  `Get passwword if exist from your vault`,
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
		prompt := &survey.Password{Message: "Enter your  master password:"}
		survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))

		vaultPwd := argon2.IDKey([]byte(masterPassword), []byte(sessionKey), 1, 64*1024, 4, 32)

		vaultData, err := crypt.DecryptFile(path, vaultPwd)

		if err != nil {
			return err
		}

		var s storage.Storage

		err = json.Unmarshal(vaultData, &s)
		if err != nil {
			return err
		}

		var name string

		namePrompt := &survey.Input{Message: "Enter name for which you want to get password:"}

		err = survey.AskOne(namePrompt, &name)
		if err != nil {
			return err
		}

		if _, ok := s.Passwords[name]; !ok {
			return errors.New("failed to find this password")
		}

		fmt.Println(string(s.Passwords[name]))

		return nil
	},
}
