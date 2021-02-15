package random_test

import (
	. "github.com/Nnachevvv/passmag/random"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Random", func() {
	Describe("Generate random string", func() {
		Context("Generate random string", func() {
			It("should return an empty string", func() {
				Expect(StringRune(0)).To(Equal(""))
			})
		})

		Context("Generate random string with length one", func() {
			It("should contains only letters from rune", func() {
				for i := 0; i < 100; i++ {
					Expect(isLetter(StringRune(1))).To(Equal(true))

				}
			})
		})
	})

})

func isLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
