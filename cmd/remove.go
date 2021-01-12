package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/nnachevv/passmag/user"
	"github.com/spf13/cobra"
)

var remove = &cobra.Command{

	Use:   "remove",
	Short: "Remove password from your password manager",
	Long:  `Remove password from your password manager`,
	RunE: func(cmd *cobra.Command, args []string) error {
		u, err := user.EnterSession()
		if err != nil {
			return err
		}

		err = removePassword(u)
		if err != nil {
			return err
		}

		fmt.Println("successfully removed password")

		return nil
	},
}

func removePassword(u user.User) error {
	var removeName string
	prompt := &survey.Password{Message: "Enter for which URL you want to remove password:"}
	survey.AskOne(prompt, &removeName, survey.WithValidator(survey.Required))

	err := survey.AskOne(prompt, &removeName)
	if err != nil {
		return err
	}

	s, err := storage.Load(u.VaultData)
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

	err = crypt.EncryptFile(u.VaultPath, byteData, u.VaultPwd())

	if err != nil {
		return fmt.Errorf("failed to encrypt sessionData : %w", err)
	}

	err = SyncVault(s, u.Password)
	if err != nil && err != ErrCreateUser {
		return err
	}

	return nil
}
