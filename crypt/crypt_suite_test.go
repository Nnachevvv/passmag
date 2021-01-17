package crypt_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCrypt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crypt Suite")
}
