package cmd

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
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
		if !viper.IsSet("PASS_KEY") {
			prompt := &survey.Input{Message: "Please enter your session key :"}
			survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required))
		} else {
			sessionKey = viper.GetString("PASS_KEY")
		}

		var masterPassword string
		prompt := &survey.Password{Message: "Enter your  master password:"}
		survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))

		var email []byte
		//storage.email()
		vaultPwd := argon2.IDKey([]byte(masterPassword), []byte(email), 1, 64*1024, 4, 32)

		//TODO :ssh
		encryptHost := argon2.IDKey([]byte(vaultPwd), []byte(answers.Host), 1, 64*1024, 4, 32)

		encryptPwd := argon2.IDKey([]byte(vaultPwd), []byte(answers.Password), 1, 64*1024, 4, 32)

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"email": vaultPwd},
			bson.D{
				{"$set", bson.D{{string(encryptHost), encryptPwd}}},
			},
		)
		if err != nil {
			return err
		}

		return nil
	},
}
