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
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/mocks"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

var _ = Describe("Remove", func() {
	var (
		c         *expect.Console
		state     *vt10x.State
		err       error
		path      string
		removeCmd *cobra.Command
		stdOut    bytes.Buffer
		stdErr    bytes.Buffer
		mockCtrl  *gomock.Controller
		mockDB    *mocks.MockDatabase

		vaultPwd []byte
	)

	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		cmd.Stdio = terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
		cmd.Crypt = crypt.Crypt{}

		mockCtrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockDatabase(mockCtrl)
		cmd.MongoDB.Database = mockDB

		removeCmd = cmd.NewRemoveCmd()
		removeCmd.SetArgs([]string{})
		removeCmd.SetOut(&stdOut)
		removeCmd.SetErr(&stdErr)

		vaultPwd = argon2.IDKey([]byte("test-dummy"), []byte("MRfbladUgDxLHvVWbxUjQUiZQykqiNcK"), 1, 64*1024, 4, 32)
		path, err = tempFile("fixtures/vault.bin")
		Expect(err).ShouldNot(HaveOccurred())

		viper.Set("password.path", path)
		viper.Set("PASS_SESSION", "MRfbladUgDxLHvVWbxUjQUiZQykqiNcK")
	})
	AfterEach(func() {

		mockCtrl.Finish()
	})

	Context("pass password which want to remove", func() {
		It("removes password from vault", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name of password you want to remove:")
				c.SendLine("exist@mail.com")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
				c.ExpectEOF()
			}()
			mockDB.EXPECT().Insert("exist@mail.com", gomock.Any())
			err = removeCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			password, ok := getAddedPassword(path, "exist@mail.com", vaultPwd)
			Expect(ok).Should(BeFalse())
			Expect(string(password)).To(Equal(""))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("pass non existing ", func() {
		It("give non existing err", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name of password you want to remove:")
				c.SendLine("nonexist@mail.com")
				c.ExpectEOF()
			}()
			err = removeCmd.Execute()
			Expect(err).To(Equal(errors.New("this name not exist in your vault")))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("pass wrong master password", func() {
		It("throw failed to find this name error", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("wrong")
				c.ExpectEOF()

			}()
			err = removeCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})
