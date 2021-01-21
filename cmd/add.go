package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nnachevv/passmag/random"

	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/nnachevv/passmag/user"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// NewAddCmd creates a new addCmd
func NewAddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Initialize email, password and master password for your password manager",
		Long:  `Set master password`,
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := user.EnterSession()
			if err != nil {
				return err
			}

			err = addPassword(u)
			if err != nil {
				return err
			}
			fmt.Println("successfully added")
			return nil
		},
	}
	return addCmd
}

func addPassword(u user.User) error {

	var name string

	namePrompt := &survey.Input{Message: "Enter name for your password:"}
	survey.AskOne(namePrompt, &name, survey.WithValidator(survey.Required))

	password, err := processPassword()
	if err != nil {
		return err
	}

	s, err := storage.Load(u.VaultData)
	if err != nil {
		return err
	}

	err = s.Add(name, password)
	if err != nil {
		var confirm bool
		editConfirm := &survey.Confirm{Message: "This name with password already exist! Do you want to edit name with newly password"}
		survey.AskOne(editConfirm, &confirm, survey.WithValidator(survey.Required))
		if !confirm {
			return nil
		}
		s.Edit(name, password)
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

	return nil
}

func processPassword() (string, error) {
	var confirm bool
	generateConfirm := &survey.Confirm{Message: "Do you want to automatically generate password?"}
	err := survey.AskOne(generateConfirm, &confirm, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("failed to get input : %w", err)
	}

	if confirm {
		return random.StringRune(32), nil
	}

	var password string
	passwordPrompt := &survey.Password{Message: "Enter your password:"}
	err = survey.AskOne(passwordPrompt, &password, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("failed to get input : %w", err)
	}

	return password, nil
}
