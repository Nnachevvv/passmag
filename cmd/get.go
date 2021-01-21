package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/storage"
	"github.com/nnachevv/passmag/user"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get password from your vault",
	Long:  `Get password if exist from your vault`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pass, err := getPassword()
		if err != nil {
			return err
		}

		fmt.Println(pass)
		return nil
	},
}

func getPassword() (string, error) {
	u, err := user.EnterSession()
	if err != nil {
		return "", err
	}

	var s storage.Storage

	err = json.Unmarshal(u.VaultData, &s)
	if err != nil {
		return "", err
	}

	var name string

	namePrompt := &survey.Input{Message: "Enter name for which you want to get password:"}

	err = survey.AskOne(namePrompt, &name)
	if err != nil {
		return "", err
	}

	pass, err := s.Get(name)
	if err != nil {
		return "", err
	}

	return pass, nil
}
