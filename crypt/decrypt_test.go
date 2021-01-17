package crypt_test

import (
	"crypto/aes"
	"errors"

	. "github.com/nnachevv/passmag/crypt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decrypts", func() {
	Describe("Decrypts already encrypted data", func() {
		Context("decrypt data with password with length different than 16,24,32", func() {
			It("should return an invalid key size error", func() {
				_, err := Decrypt([]byte("dummy-value"), []byte("dummy-password"))
				expectedError := aes.KeySizeError(len("dummy-password"))
				Expect(err).To(Equal(expectedError))
			})
		})

		Context("decrypt data with wrong password", func() {
			It("should return an invalid key size error", func() {
				encryptedBytes := []byte{91, 143, 41, 169, 177, 215, 127, 148, 100, 39, 80, 254, 182, 55, 238, 64, 70, 68, 4, 112, 79, 145, 142, 188, 135, 161, 143, 166, 63}
				_, err := Decrypt(encryptedBytes, []byte("11111111111111111111111111111111"))
				expectedErr := errors.New("cipher: message authentication failed")
				Expect(err).To(Equal(expectedErr))
			})
		})

		Context("decrypt password with valid password", func() {
			It("should return value of string which is one", func() {
				encryptedBytes := []byte{90, 143, 41, 169, 177, 215, 127, 148, 100, 39, 80, 254, 182, 55, 238, 64, 70, 68, 4, 112, 79, 145, 142, 188, 135, 161, 143, 166, 63}
				val, err := Decrypt(encryptedBytes, []byte("11111111111111111111111111111111"))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(val)).To(Equal("1"))
			})
		})

	})
})
