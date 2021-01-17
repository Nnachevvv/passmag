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

var edit = &cobra.Command{
	Use:   "edit",
	Short: "Set new name for password",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		u, err := user.EnterSession()
		if err != nil {
			return err
		}

		var s storage.Storage

		err = json.Unmarshal(u.VaultData, &s)
		if err != nil {
			return err
		}

		editQs := []*survey.Question{
			{
				Name:   "name",
				Prompt: &survey.Input{Message: "Enter name for which you want to change password:"},
			},
			{
				Name:   "newname",
				Prompt: &survey.Input{Message: "Enter new name for your password:"},
			},
		}

		answers := struct {
			Name    string
			NewName string
		}{}

		err = survey.Ask(editQs, &answers)
		if err != nil {
			return fmt.Errorf("failed to get input : %w", err)
		}

		pwd, err := s.Get(answers.Name)
		if err != nil {
			return err
		}

		err = s.Remove(answers.Name)
		if err != nil {
			return err
		}

		err = s.Add(answers.NewName, pwd)
		if err != nil {
			return err
		}

		byteData, err := json.Marshal(s)
		if err != nil {
			return fmt.Errorf("failed to marshal map : %w", err)
		}
		err = crypt.EncryptFile(u.VaultPath, byteData, u.VaultPwd)

		if err != nil {
			return fmt.Errorf("failed to encrypt sessionData : %w", err)
		}

		err = SyncVault(s, u.Password)
		if err != nil && err != ErrCreateUser {
			return err
		}
		fmt.Println("succesfuly moved your password")

		return nil
	},
}
