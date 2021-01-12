package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nnachevv/passmag/cmd/mongo"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(login)
	rootCmd.AddCommand(initialize)
	rootCmd.AddCommand(remove)
	rootCmd.AddCommand(add)
	rootCmd.AddCommand(get)
	rootCmd.AddCommand(edit)
	rootCmd.AddCommand(change)
	rootCmd.AddCommand(logout)

	//defer client.Disconnect(ctx)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
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
		return fmt.Errorf("failed to add user to db", err)
	}

	service.Insert(s.Email, vaultData)

	if err != nil {
		return ErrCreateUser
	}

	return nil
}
