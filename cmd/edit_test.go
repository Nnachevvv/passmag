package cmd_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/Nnachevvv/passmag/cmd"
	"github.com/Nnachevvv/passmag/crypt"
	"github.com/Nnachevvv/passmag/mocks"
	"github.com/golang/mock/gomock"
	"github.com/hinshun/vt10x"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

var _ = Describe("Edit", func() {
	var (
		c            *expect.Console
		state        *vt10x.State
		err          error
		path         string
		editCmd      *cobra.Command
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
		cmd.MongoDB = mockDatabase
		editCmd = cmd.NewEditCmd()
		editCmd.SetArgs([]string{})
		editCmd.SetOut(&stdOut)
		editCmd.SetErr(&stdErr)

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
		It("gets vault, and change name of password", func() {
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter existing name in your vault:")
				c.SendLine("exist@mail.com")
				c.ExpectString("Enter new name for your password:")
				c.SendLine("new@mail.com")
				c.ExpectEOF()
			}()
			mockDatabase.EXPECT().Insert("exist@mail.com", gomock.Any(), gomock.Any())
			err = editCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			password, ok := getAddedPassword(path, "exist@mail.com", vaultPwd)

			Expect(ok).Should(BeFalse())
			Expect(string(password)).To(Equal(""))

			password, ok = getAddedPassword(path, "new@mail.com", vaultPwd)
			Expect(ok).Should(BeTrue())
			Expect(string(password)).To(Equal("gMdLasZIGAEmDSCprqFkZQSAnjzeZzUP"))

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
				c.ExpectString("Enter existing name in your vault:")
				c.SendLine("nonexist@mail.com")
				c.ExpectString("Enter new name for your password:")
				c.SendLine("new@mail.com")
				c.ExpectEOF()

			}()

			err = editCmd.Execute()
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
			err = editCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})
