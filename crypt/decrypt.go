package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io/ioutil"
)

//Decrypt decrypts data by given passphrase
func (c Crypt) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
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

// DecryptFile decrypts file by given password
func (c Crypt) DecryptFile(filename string, key []byte) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	byteEncrypted, err := c.Decrypt(data, key)
	if err != nil {
		return []byte{}, err
	}

	return byteEncrypted, nil
}
