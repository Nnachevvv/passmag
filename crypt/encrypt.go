package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
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

// EncryptFile encrypts given data with cypher algorithm and saves it to file.
func (c Crypt) EncryptFile(filename string, data []byte, key []byte) error {
	pathDir := filepath.Dir(filename)

	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		err := os.MkdirAll(pathDir, 0777)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	byteEncrypted, err := Encrypt(data, key)
	if err != nil {
		return err
	}

	f.Write(byteEncrypted)
	return nil
}
