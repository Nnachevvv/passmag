package cmd_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/golang/mock/gomock"
	"github.com/hinshun/vt10x"
	"github.com/nnachevv/passmag/cmd"
	"github.com/nnachevv/passmag/mocks"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
)

var _ = Describe("Login", func() {
	var (
		c        *expect.Console
		state    *vt10x.State
		err      error
		loginCmd *cobra.Command
		stdOut   bytes.Buffer
		stdErr   bytes.Buffer

		mockCtrl     *gomock.Controller
		mockDatabase *mocks.MockDatabase
		mockCrypt    *mocks.MockCrypter
	)

	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		cmd.Stdio = terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}

		mockCtrl = gomock.NewController(GinkgoT())

		mockDatabase = mocks.NewMockDatabase(mockCtrl)
		mockCrypt = mocks.NewMockCrypter(mockCtrl)
		cmd.MongoDB.Database = mockDatabase
		loginCmd = cmd.NewLoginCmd()
		cmd.Crypt = mockCrypt

		loginCmd.SetArgs([]string{})
		loginCmd.SetOut(&stdOut)
		loginCmd.SetErr(&stdErr)
	},
	)

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("with valid account", func() {
		It("contains account in db", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your email address:")
				c.SendLine("dummy")
				c.ExpectString("email should be longer than 8 characters")
				c.SendLine("test-dummy2@mail.com")
				c.ExpectString("Enter your  master password:")
				c.SendLine("test")
				c.ExpectString("password should be longer than 8 characters")
				c.SendLine("test-dummy")
				c.ExpectEOF()
			}()
			expectedEmail := "dummytest-dummy2@mail.com"
			expectedPassword := "test-dummy"
			vaultPwd := argon2.IDKey([]byte(expectedPassword), []byte(expectedEmail), 1, 64*1024, 4, 32)

			mockDatabase.EXPECT().Find("dummytest-dummy2@mail.com").Return(primitive.M{"vault": primitive.Binary{Data: []byte("testValueIn")}}, nil)
			mockCrypt.EXPECT().Decrypt(gomock.Any(), vaultPwd)
			mockCrypt.EXPECT().EncryptFile(gomock.Any(), gomock.Any(), gomock.Any())

			err = loginCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("with wrong password", func() {
		It("returns err", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your email address:")
				c.SendLine("dummy")
				c.ExpectString("email should be longer than 8 characters")
				c.SendLine("test-dummy2@mail.com")
				c.ExpectString("Enter your  master password:")
				c.SendLine("test")
				c.ExpectString("password should be longer than 8 characters")
				c.SendLine("test-dummy")
				c.ExpectEOF()
			}()
			expectedEmail := "dummytest-dummy2@mail.com"
			expectedPassword := "test-dummy"
			vaultPwd := argon2.IDKey([]byte(expectedPassword), []byte(expectedEmail), 1, 64*1024, 4, 32)
			err := errors.New("Mock Error")
			mockDatabase.EXPECT().Find("dummytest-dummy2@mail.com").Return(primitive.M{"vault": primitive.Binary{Data: []byte("testValueIn")}}, nil)
			mockCrypt.EXPECT().Decrypt(gomock.Any(), vaultPwd).Return([]byte{}, err)

			err = loginCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("with account that not exist", func() {
		It("returns err", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your email address:")
				c.SendLine("dummy")
				c.ExpectString("email should be longer than 8 characters")
				c.SendLine("test-dummy2@mail.com")
				c.ExpectString("Enter your  master password:")
				c.SendLine("test")
				c.ExpectString("password should be longer than 8 characters")
				c.SendLine("test-dummy")
				c.ExpectEOF()
			}()

			err := errors.New("Mock Error")
			mockDatabase.EXPECT().Find("dummytest-dummy2@mail.com").Return(primitive.M{}, err)

			err = loginCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

})
