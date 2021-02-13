package cmd

import (
	"fmt"
	"os"

	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
)

// NewLogoutCmd creates a new logoutCmd
func NewLogoutCmd() *cobra.Command {
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from logged user",
		Long:  `Logout from logged user and delete currently download vault`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := storage.FilePath()
			if err != nil {
				return err
			}

			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to remove vault from %s, please delete it manually: %w", path, err)
			}

			return nil
		},
	}
	return logoutCmd
}
