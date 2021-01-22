package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nnachevv/passmag/cmd/mongo"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

var (
	// Used for flags.
	cfgFile string
	service mongo.Service

	rootCmd = &cobra.Command{
		Use:   "passmag",
		Short: "A password manager used to store securely passwords",
		Long:  `passmag`,
	}
)

func init() {
	io := terminal.Stdio{os.Stdin, os.Stdout, os.Stderr}

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(initializeCmd)
	rootCmd.AddCommand(NewRemoveCmd(io))
	rootCmd.AddCommand(NewAddCmd(io))
	rootCmd.AddCommand(NewGetCmd(io))
	rootCmd.AddCommand(NewCopyCmd(io))
	rootCmd.AddCommand(NewEditCmd(io))
	rootCmd.AddCommand(NewChangeCmd(io))
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(NewListCmd(io))

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

// ErrCreateUser throw by db when try to insert user
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
