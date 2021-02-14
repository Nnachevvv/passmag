package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new listCmd
func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all password from your vault",
		Long:  "Lists all password from your vault in plain text",
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

			for n, p := range s.Passwords {
				fmt.Fprintf(cmd.OutOrStdout(), "%s : %s\n", n, string(p))
			}

			return nil
		},
	}
	return listCmd
}
