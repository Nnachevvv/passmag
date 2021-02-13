package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nnachevv/passmag/cmd/mongo"
	"github.com/nnachevv/passmag/crypt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	service mongo.Service
	// Stdio is used for testing virtuall terminal.
	Stdio   terminal.Stdio
	rootCmd = &cobra.Command{
		Use:   "passmag",
		Short: "A password manager used to store securely passwords",
		Long:  `passmag`,
	}

	// Crypt is used to mock and encapsulated abstract encrypt functionality.
	Crypt crypt.Crypter

	// MongoDB is used to mock and encapsulated abstract database functionality.
	MongoDB mongo.MongoDatabase
)

func init() {
	Stdio = terminal.Stdio{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}

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

	service.Connect()
	viper.AutomaticEnv()
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

/*// ErrCreateUser throw by db when try to insert user
var ErrCreateUser = errors.New("failed to add user to db")

// SyncVault syncs current state of vault to password if internet connection is provided
func SyncVault(s storage.Storage, password []byte) error {
	byteData, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal map : %w", err)
	}

	vaultPwd := argon2.IDKey(password, []byte(s.Email), 1, 64*1024, 4, 32)
	vaultData, err := crypt.Encrypt(byteData, vaultPwd)
	if err != nil {
		return fmt.Errorf("failed to add user to db :%w", err)
	}

	s.TimeCreated = time.Now()
	service.Insert(s.Email, vaultData)

	if err != nil {
		return ErrCreateUser
	}

	return nil
}
*/
