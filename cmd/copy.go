package cmd

import (
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{

	Use:   "cp",
	Short: "Copy password to cpliboard",
	Long:  `Get password if exist and copy to clipboard.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pass, err := getPassword()
		if err != nil {
			return err
		}
		clipboard.WriteAll(pass)
		return nil
	},
}
