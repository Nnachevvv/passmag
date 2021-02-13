package crypt

//Crypter interface abstacts encrypting and decrypting values
type Crypter interface {
	EncryptFile(string, []byte, []byte) error
	DecryptFile(string, []byte) ([]byte, error)
	Decrypt([]byte, []byte) ([]byte, error)
}

//Crypt is struct used to represent crypting functions
type Crypt struct{}
