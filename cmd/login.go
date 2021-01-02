package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nnachevv/passmag/crypt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/argon2"
)

// the questions to ask
var loginQs = []*survey.Question{
	{
		Name:   "email",
		Prompt: &survey.Input{Message: "Enter your email adress:"},
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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to password manager CLI",
	Long:  "login to password manager CLI and seal vault locally with generated random decription key",
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := struct {
			Email          string
			MasterPassword string
		}{}

		db := struct {
			vaultPwd string `bson:"vaultPwd"`
		}{}

		err := survey.Ask(loginQs, &answers)
		if err != nil {
			return err
		}
		fmt.Println(answers.MasterPassword)
		fmt.Println(answers.Email)

		vaultPwd := argon2.IDKey([]byte(answers.MasterPassword), []byte(answers.Email), 1, 64*1024, 4, 32)

		collection.FindOne(context.Background(), bson.M{"vaultPwd": vaultPwd}).Decode(&db)

		fmt.Println(string(vaultPwd))
		fmt.Println("---------------------")
		fmt.Println(string(db.vaultPwd))
		fmt.Println("---------------------")

		if db.vaultPwd == "" || string(vaultPwd) != db.vaultPwd {
			return errors.New("failed to find this account")
		}

		var record bson.M
		collection.FindOne(context.Background(), bson.M{"vaultPwd": vaultPwd}).Decode(&record)

		s := storage.New(record, answers.Email)

		json, err := json.Marshal(s)
		if err != nil {
			return fmt.Errorf("failed to marshal storage into struct : %w", err)
		}

		fmt.Println(json)
		sessionKey := make([]byte, 20)
		fmt.Println("sessionKey")
		crypt.EncryptFile("test.bin", json, string(sessionKey))

		test := crypt.DecryptFile("test.bin", string(sessionKey))
		fmt.Println(string(test))

		//sync data
		/*b, err := json.Marshal(&record)
		if err != nil {
			panic(err) // it will be invoked
			// panic: json: unsupported value: NaN
		}

		f, err := os.Create("data.txt")

		f.WriteString(string(b))


		fmt.Printf("Your sessionKey is %s", string(sessionKey))

		//encrypt data
		*/
		return nil
	},
}

func setSessionKey() {

}

/*
func encryptFile(sessionKey []byte, masterPassword []byte, vault storage) {

}
*/
