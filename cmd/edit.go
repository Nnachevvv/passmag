package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
)

// the questions to ask
var editQs = []*survey.Question{
	{
		Name:   "name",
		Prompt: &survey.Input{Message: "Enter name for which you want to change password:"},
	},
	{
		Name:   "newname",
		Prompt: &survey.Input{Message: "Enter new name for your password:"},
	},
}

var edit = &cobra.Command{
	Use:   "edit",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultData, vaultPwd, path, err := EnterSession()
		if err != nil {
			return err
		}

		var s storage.Storage

		err = json.Unmarshal(vaultData, &s)
		if err != nil {
			return err
		}

		answers := struct {
			Name    string
			NewName string
		}{}

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
		err = crypt.EncryptFile(path, byteData, vaultPwd)

		if err != nil {
			return fmt.Errorf("failed to encrypt sessionData : %w", err)
		}
		return nil
		fmt.Println("succesfuly moved your password")

		return nil
	},
}
