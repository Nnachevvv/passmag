package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
)

// NewGetCmd creates a new getCmd
func NewGetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get password from your vault",
		Long:  `Get password if exist , otherwise error will be thrown`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pass, err := getPassword()
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), pass)
			return nil
		},
	}
	return getCmd
}

func getPassword() (string, error) {
	u, err := EnterSession()
	if err != nil {
		return "", err
	}

	var s storage.Storage

	err = json.Unmarshal(u.VaultData, &s)
	if err != nil {
		return "", err
	}

	var name string

	namePrompt := &survey.Input{Message: "Enter name for which you want to get your password:"}

	err = survey.AskOne(namePrompt, &name, survey.WithStdio(Stdio.In, Stdio.Out, Stdio.Err))
	if err != nil {
		return "", err
	}

	pass, err := s.Get(name)
	if err != nil {
		return "", err
	}

	return pass, nil
}
