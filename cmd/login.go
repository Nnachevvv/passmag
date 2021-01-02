package cmd

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nnachevv/passmag/crypt"

	"github.com/AlecAivazis/survey/v2"
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

		err := survey.Ask(loginQs, &answers)
		if err != nil {
			return err
		}

		err = CheckAuthorization(answers.Email, answers.MasterPassword)
		if err != nil {
			return err
		}
		var record bson.M
		collection.FindOne(context.Background(), bson.M{"email": answers.Email}).Decode(&record)

		//s := storage.New(record, answers.Email)

		//json, err := json.Marshal(s)
		//if err != nil {
		//	return fmt.Errorf("failed to marshal storage into struct : %w", err)
		//}

		//sessionKey := make([]byte, 20)
		//crypt.EncryptFile("test.bin", json, string(sessionKey))

		//test := crypt.DecryptFile("test.bin", string(sessionKey))
		//fmt.Println(string(test))

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

func CheckAuthorization(email string, password string) error {
	db := struct {
		Email string `bson:"email"`
	}{}

	collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&db)
	fmt.Printf("e: %s", db.Email)
	if db.Email == "" || string(email) != db.Email {
		return errors.New("failed to find this account")
	}

	var record bson.M
	collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&record)
	vaultPwd := argon2.IDKey([]byte(password), []byte(email), 1, 64*1024, 4, 32)

	_, err := crypt.Decrypt(record["vault"].(primitive.Binary).Data, vaultPwd)
	if err != nil {
		return fmt.Errorf("failed to decrypt data value: %w ", err)
	}
	return nil

}
