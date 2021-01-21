package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

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
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(initializeCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(NewAddCmd())
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(changeCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(listCmd)

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

	vaultPwd := argon2.IDKey([]byte(s.Email), password, 1, 64*1024, 4, 32)
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
