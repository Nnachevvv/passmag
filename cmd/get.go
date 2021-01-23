package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nnachevv/passmag/storage"
	"github.com/nnachevv/passmag/user"
	"github.com/spf13/cobra"
)

// NewGetCmd creates a new getCmd
func NewGetCmd(stdio terminal.Stdio) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get password from your vault",
		Long:  `Get password if exist from your vault`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pass, err := getPassword(stdio)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), pass)
			return nil
		},
	}
	return getCmd
}

func getPassword(stdio terminal.Stdio) (string, error) {
	u, err := user.EnterSession(stdio)
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

	err = survey.AskOne(namePrompt, &name, survey.WithStdio(stdio.In, stdio.Out, stdio.Err))
	if err != nil {
		return "", err
	}

	pass, err := s.Get(name)
	if err != nil {
		return "", err
	}

	return pass, nil
}
