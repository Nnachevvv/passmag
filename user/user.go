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
	Password   []byte
	SessionKey []byte
	VaultPath  string
	VaultData  []byte
}

// VaultPwd returns vault password for given user
func (u *User) VaultPwd() []byte {
	vaultPwd := argon2.IDKey(u.Password, u.SessionKey, 1, 64*1024, 4, 32)
	return vaultPwd
}

// LoadVault load current vualt from user directory
func (u *User) LoadVault() error {
	vaultData, err := crypt.DecryptFile(u.VaultPath, u.VaultPwd())
	if err != nil {
		return fmt.Errorf("failed to load your vault try again : %w", err)
	}

	u.VaultData = vaultData
	return nil
}

//EnterSession prompts to enter Session Key and ask for master password
func EnterSession() (User, error) {
	path, err := storage.FilePath()
	if err != nil {
		return User{}, err
	}

	if err := storage.VaultExist(path); err != nil {
		return User{}, err
	}

	var sessionKey string
	if !viper.IsSet("PASS_SESSION") {
		prompt := &survey.Input{Message: "Please enter your session key :"}
		survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required))
	} else {
		sessionKey = viper.GetString("PASS_SESSION")
	}

	var masterPassword string
	prompt := &survey.Password{Message: "Enter your  master password:"}
	survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))
	u := User{
		Password:   []byte(masterPassword),
		SessionKey: []byte(sessionKey),
		VaultPath:  path}

	err = u.LoadVault()
	if err != nil {
		return User{}, err
	}

	return u, err
}
