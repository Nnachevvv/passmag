package crypt_test

import (
	"crypto/aes"

	. "github.com/Nnachevvv/passmag/crypt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encrypt", func() {
	Describe("Encrypts data with given value", func() {
		Context("encrypt data with password with length different than 16,24,32", func() {
			It("should return an invalid key size error", func() {
				_, err := Encrypt([]byte("dummy-value"), []byte("dummy-password"))
				expectedError := aes.KeySizeError(len("dummy-password"))
				Expect(err).To(Equal(expectedError))
			})
		})

		Context("encrypt password with key with size 32", func() {
			It("should not return error", func() {
				_, err := Encrypt([]byte("1"), []byte("11111111111111111111111111111111"))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
