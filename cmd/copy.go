package cmd

import (
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// NewCopyCmd creates a new copyCmd
func NewCopyCmd() *cobra.Command {
	copyCmd := &cobra.Command{
		Use:   "cp",
		Short: "Copy password to cpliboard",
		Long:  `Copy password to clipboard if exist , otherwise error will be thrown.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pass, err := getPassword()
			if err != nil {
				return err
			}

			clipboard.WriteAll(pass)
			return nil
		},
	}
	return copyCmd
}
