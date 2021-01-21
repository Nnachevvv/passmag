package user

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

// User contains current logged user
type User struct {
	Password  []byte
	VaultPwd  []byte
	VaultPath string
	VaultData []byte
}

//EnterSession prompts to enter Session Key and ask for master password
func EnterSession() (User, error) {
	var path string
	var err error
	if !viper.IsSet("path") {
		path, err = storage.FilePath()
		if err != nil {
			return User{}, err
		}
	} else {
		path = viper.GetString("path")
	}

	if err := storage.VaultExist(path); err != nil {
		return User{}, err
	}

	var sessionKey, masterPassword string
	if !viper.IsSet("PASS_SESSION") {
		prompt := &survey.Input{Message: "Please enter your session key :"}
		survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required))
	} else {
		sessionKey = viper.GetString("PASS_SESSION")
	}

	prompt := &survey.Password{Message: "Enter your  master password:"}
	survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))
	u := User{
		Password:  []byte(masterPassword),
		VaultPwd:  argon2.IDKey([]byte(masterPassword), []byte(sessionKey), 1, 64*1024, 4, 32),
		VaultPath: path}

	err = u.loadVault()
	if err != nil {
		return User{}, err
	}

	return u, err
}

func (u *User) loadVault() error {
	vaultData, err := crypt.DecryptFile(u.VaultPath, u.VaultPwd)
	if err != nil {
		return fmt.Errorf("failed to load your vault try again : %w", err)
	}

	u.VaultData = vaultData
	return nil
}
