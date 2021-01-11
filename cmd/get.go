package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
)

var get = &cobra.Command{

	Use:   "get",
	Short: "Get password from your vault",
	Long:  `Get passwword if exist from your vault`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultData, _, _, err := EnterSession()
		if err != nil {
			return err
		}

		var s storage.Storage

		err = json.Unmarshal(vaultData, &s)
		if err != nil {
			return err
		}

		var name string

		namePrompt := &survey.Input{Message: "Enter name for which you want to get password:"}

		err = survey.AskOne(namePrompt, &name)
		if err != nil {
			return err
		}

		if _, ok := s.Passwords[name]; !ok {
			return errors.New("failed to find this password")
		}

		fmt.Println(string(s.Passwords[name]))

		return nil
	},
}
