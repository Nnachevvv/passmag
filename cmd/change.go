package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/nnachevv/passmag/user"
	"github.com/spf13/cobra"
)

var change = &cobra.Command{

	Use:   "change",
	Short: "Change password for given host",
	Long:  `Change password for given host`,
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

		answers := struct {
			Name     string
			Password string
		}{}

		changeQs := []*survey.Question{
			{
				Name:   "name",
				Prompt: &survey.Input{Message: "Enter name for which you want to edit your password:"},
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

		err = survey.Ask(changeQs, &answers)
		if err != nil {
			return err
		}

		if _, ok := s.Passwords[answers.Name]; !ok {
			return errors.New("failed to find this name")
		}

		s.Edit(answers.Name, answers.Password)
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

		return nil
	},
}
