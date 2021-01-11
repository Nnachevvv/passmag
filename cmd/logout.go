package cmd

import (
	"fmt"
	"os"

	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
)

var logout = &cobra.Command{

	Use:   "logout",
	Short: "Logout from logged user",
	Long:  `Logout from logged user and deleate currently downoload vault`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := storage.FilePath()
		if err != nil {
			return err
		}

		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("failed to remove vault from %s, please deleate it manually: %w", path, err)
		}

		return nil
	},
}
