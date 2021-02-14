package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
)

// NewEditCmd creates a new editCmd
func NewEditCmd() *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Set new name for password",
		Long:  `Ask name of password which user will want to change `,
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := EnterSession()
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
					Prompt: &survey.Input{Message: "Enter existing name in your vault:"},
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

			err = survey.Ask(editQs, &answers, survey.WithStdio(Stdio.In, Stdio.Out, Stdio.Err))
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
			err = Crypt.EncryptFile(u.VaultPath, byteData, u.VaultPwd)

			if err != nil {
				return fmt.Errorf("failed to encrypt sessionData : %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "successfully moved your password")
			s.SyncStorage(u.Password, MongoDB, Client)

			return nil
		},
	}
	return editCmd
}
