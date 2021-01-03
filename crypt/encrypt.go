package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

//Encrypt encrypts data by given passphrase
func Encrypt(data []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

//Decrypt decrypts data by given passphrase
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptFile encripts given data with cypher algorithm and saves it to file.
func EncryptFile(filename string, data []byte, key []byte) error {
	f, _ := os.Create(filename)
	defer f.Close()
	byteEncrypted, err := Encrypt(data, key)
	if err != nil {
		return err
	}

	f.Write(byteEncrypted)
	return nil
}

// DecryptFile decrypts file by given password
func DecryptFile(filename string, key []byte) ([]byte, error) {
	data, _ := ioutil.ReadFile(filename)

	byteEncrypted, err := Decrypt(data, key)
	if err != nil {
		return []byte{}, err
	}

	return Decrypt(byteEncrypted, key)
}
