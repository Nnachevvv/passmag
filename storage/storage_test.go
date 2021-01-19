package storage_test

import (
	"errors"
	"time"

	"github.com/nnachevv/passmag/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("storage package", func() {
	Context("storage.New()", func() {
		When("pass current time and data which contains email and dummy values", func() {
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

		When("pass data without email", func() {
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

	Context("storage.Add()", func() {
		var (
			s   storage.Storage
			err error
		)
		BeforeEach(func() {
			s.Passwords = make(map[string][]byte)
			err = s.Add("test", "dummy")

		})

		When("add name, password", func() {
			It("returns password containing this name", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(s.Passwords["test"])).To(Equal("dummy"))
			})
		})

		When("add name, password that exist in db", func() {
			It("returns error: this name exist in db", func() {
				err = s.Add("test", "dummy")
				expectedErr := errors.New("this name already exist")
				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Context("storage.Remove()", func() {
		var (
			s   storage.Storage
			err error
		)
		BeforeEach(func() {
			s.Passwords = make(map[string][]byte)

		})

		When("remove name, password", func() {
			It("returns nil and name don't exist in db", func() {
				s.Passwords["test"] = []byte("dummy")
				err := s.Remove("test")
				Expect(err).ShouldNot(HaveOccurred())
				_, ok := s.Passwords["test"]
				Expect(ok).To(BeFalse())
			})
		})

		When("add name, password that exist in db", func() {
			It("returns error: this name not exist in db", func() {
				err = s.Remove("test")
				expectedErr := errors.New("this name not exist in our db")
				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Context("storage.Get()", func() {
		var (
			s   storage.Storage
			err error
		)
		BeforeEach(func() {
			s.Passwords = make(map[string][]byte)

		})

		When("gets name", func() {
			It("returns password", func() {
				s.Passwords["test"] = []byte("dummy")
				name, err := s.Get("test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(name).To(Equal("dummy"))
			})
		})

		When("get name that not exist in db", func() {
			It("returns error: this name not exist in your password manager", func() {
				_, err = s.Get("test1")
				expectedErr := errors.New("this name not exist in your password manager ")
				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Context("storage.Edit()", func() {
		var (
			s storage.Storage
		)
		BeforeEach(func() {
			s.Passwords = make(map[string][]byte)
			s.Passwords["test"] = []byte("dummy")
		})

		When("edits name", func() {
			It("returns edited password", func() {
				s.Edit("test", "dummy2")
				Expect(s.Passwords["test"]).Should(Equal([]byte("dummy2")))
			})
		})
	})

})
