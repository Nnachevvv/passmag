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

var _ = Describe("Change", func() {
	var (
		c            *expect.Console
		state        *vt10x.State
		err          error
		path         string
		changeCmd    *cobra.Command
		stdOut       bytes.Buffer
		stdErr       bytes.Buffer
		vaultPwd     []byte
		mockCtrl     *gomock.Controller
		mockDatabase *mocks.MockDatabase
	)

	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		cmd.Stdio = terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
		cmd.Crypt = crypt.Crypt{}

		mockCtrl = gomock.NewController(GinkgoT())
		mockDatabase = mocks.NewMockDatabase(mockCtrl)
		cmd.MongoDB.Database = mockDatabase
		changeCmd = cmd.NewChangeCmd()
		changeCmd.SetArgs([]string{})
		changeCmd.SetOut(&stdOut)
		changeCmd.SetErr(&stdErr)

		vaultPwd = argon2.IDKey([]byte("test-dummy"), []byte("MRfbladUgDxLHvVWbxUjQUiZQykqiNcK"), 1, 64*1024, 4, 32)

		path, err = tempFile("fixtures/vault.bin")
		Expect(err).ShouldNot(HaveOccurred())

		viper.Set("password.path", path)
		viper.Set("PASS_SESSION", "MRfbladUgDxLHvVWbxUjQUiZQykqiNcK")
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("pass existing valid name", func() {
		It("gets vault, and change password for given name", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for which you want to change your password:")
				c.SendLine("exist@mail.com")
				c.ExpectString("Enter new password:")
				c.SendLine("new-password")
				c.ExpectEOF()
			}()
			mockDatabase.EXPECT().Insert("exist@mail.com", gomock.Any())
			err = changeCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			password, ok := getAddedPassword(path, "exist@mail.com", vaultPwd)

			Expect(ok).Should(BeTrue())
			Expect(string(password)).To(Equal("new-password"))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("pass not existing name", func() {
		It("throw failed to find this name error", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for which you want to change your password:")
				c.SendLine("dummy-random")
				c.ExpectString("Enter new password:")
				c.SendLine("dummy-random")
				c.ExpectEOF()

			}()
			err = changeCmd.Execute()
			Expect(err).To(Equal(errors.New("failed to find this name")))

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
			err = changeCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})
