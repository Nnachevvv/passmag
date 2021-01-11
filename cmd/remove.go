package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

var remove = &cobra.Command{

	Use:   "remove",
	Short: "Remove password from your password manager",
	Long:  `Remove password from your password manager`,
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

		err = removePassword(vaultData, path, vaultPwd)
		if err != nil {
			return err
		}

		fmt.Println("succesfully removed password")

		return nil
	},
}

func removePassword(vaultData []byte, path string, vaultPwd []byte) error {
	var removeName string
	prompt := &survey.Password{Message: "Enter for which URL you want to remove password:"}
	survey.AskOne(prompt, &removeName, survey.WithValidator(survey.Required))

	err := survey.AskOne(prompt, &removeName)
	if err != nil {
		return err
	}

	s, err := storage.Load(vaultData)
	if err != nil {
		return err
	}

	err = s.Remove(removeName)
	if err != nil {
		return err
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
