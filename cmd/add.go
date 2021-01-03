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
		Prompt: &survey.Input{Message: "Enter your password:"},
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
		answers := struct {
			Host     string
			Password string
		}{}

		err := survey.Ask(addQs, &answers)
		if err != nil {
			return err
		}

		/*TODO: check session key
		ask for master password
		host could be only hashed
		added data (password) must be encrypted not hashed, beacuse if it's hashed we cant get it back !

		*/

		var sessionKey string
		if !viper.IsSet("PASS_SESSION") {
			prompt := &survey.Input{Message: "Please enter your session key :"}
			survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required))
		} else {
			sessionKey = viper.GetString("PASS_KEY")
		}

		var masterPassword string
		prompt := &survey.Password{Message: "Enter your  master password:"}
		survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))

		vaultPwd := argon2.IDKey([]byte(masterPassword), []byte(sessionKey), 1, 64*1024, 4, 32)
		path := storage.OperatingSystem() + "vault.bin"
		vaultData, err := crypt.DecryptFile(path, vaultPwd)
		if err != nil {
			return err
		}
		var jsonVault map[string]interface{}
		jsonVault[answers.Host] = answers.Password

		err = json.Unmarshal(vaultData, &jsonVault)
		if err != nil {
			return err
		}

		s := storage.New(jsonVault)
		fmt.Println(s)
		fmt.Println(jsonVault)
		return nil
	},
}
