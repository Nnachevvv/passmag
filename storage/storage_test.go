package storage_test

import (
	"errors"
	"time"

	"github.com/nnachevv/passmag/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("New Storage", func() {
	Context("pass current time and data which contains email and dummy values", func() {
		It("contains current time and passwords should contains only dummy values ,without email", func() {
			currentTime := time.Now()
			data := map[string]interface{}{
				"email":  "dummy-email",
				"dummy1": primitive.Binary{Data: []byte("dummy-email")},
				"dummy2": primitive.Binary{Data: []byte("dummy-email1")},
				"dummy3": primitive.Binary{Data: []byte("dummy-email2")},
			}

			st, err := storage.New(data, currentTime)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(st.TimeCreated).To(Equal(currentTime))
			expectedKeys := []string{"dummy1", "dummy2", "dummy3"}
			expectedValues := [][]byte{[]byte("dummy-email"), []byte("dummy-email1"), []byte("dummy-email2")}

			for i := 0; i < len(expectedKeys); i++ {
				value, ok := st.Passwords[expectedKeys[i]]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal(expectedValues[i]))
			}

			_, ok := st.Passwords["email"]
			Expect(ok).To(Equal(false))

		})
	})

	Context("pass data without email", func() {
		It("should throw nil interface error", func() {
			currentTime := time.Now()
			data := map[string]interface{}{
				"dummy1": primitive.Binary{Data: []byte("dummy-email")},
				"dummy2": primitive.Binary{Data: []byte("dummy-email1")},
				"dummy3": primitive.Binary{Data: []byte("dummy-email2")},
			}
			_, err := storage.New(data, currentTime)
			expectedError := errors.New("failed to get entry only with passwords")
			Expect(err).To(Equal(expectedError))

		})
	})
})
