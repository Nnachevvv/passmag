package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	mongoCli "github.com/nnachevv/passmag/cmd/mongo"
	"github.com/nnachevv/passmag/crypt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	cfgFile string
	// MongoDB is used to mock and encapsulated abstract database functionality.
	MongoDB mongoCli.Database
	// Stdio is used for testing virtuall terminal.
	Stdio   terminal.Stdio
	rootCmd = &cobra.Command{
		Use:   "passmag",
		Short: "A password manager used to store securely passwords",
		Long:  `passmag`,
	}

	// Client is used to mock represent Client for mongodb
	Client *mongo.Client
	// Crypt is used to mock and encapsulated abstract encrypt functionality.
	Crypt crypt.Crypter
)

func init() {
	Stdio = terminal.Stdio{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}
	MongoDB = &mongoCli.Service{}
	Crypt = &crypt.Crypt{}
	rootCmd.AddCommand(NewLoginCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewRemoveCmd())
	rootCmd.AddCommand(NewAddCmd())
	rootCmd.AddCommand(NewGetCmd())
	rootCmd.AddCommand(NewCopyCmd())
	rootCmd.AddCommand(NewEditCmd())
	rootCmd.AddCommand(NewChangeCmd())
	rootCmd.AddCommand(NewLogoutCmd())
	rootCmd.AddCommand(NewListCmd())
	Client = mongoCli.Connect()

	viper.AutomaticEnv()
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
