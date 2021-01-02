package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// the questions to ask
var getQs = []*survey.Question{
	{
		Name:   "host",
		Prompt: &survey.Input{Message: "Enter host for which you want to get password:"},
	},
}

var getCmd = &cobra.Command{

	Use:   "get",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {

		answers := struct {
			Host     string
			Password string
		}{}

		err := survey.Ask(getQs, &answers)
		if err != nil {
			return err
		}

		return nil
	},
}
