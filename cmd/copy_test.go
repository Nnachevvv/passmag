package cmd_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/Nnachevvv/passmag/cmd"
	"github.com/Nnachevvv/passmag/crypt"
	"github.com/hinshun/vt10x"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ = Describe("Copy", func() {
	var (
		c       *expect.Console
		state   *vt10x.State
		err     error
		path    string
		copyCmd *cobra.Command
		stdOut  bytes.Buffer
		stdErr  bytes.Buffer
	)

	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		cmd.Stdio = terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
		cmd.Crypt = crypt.Crypt{}

		copyCmd = cmd.NewCopyCmd()
		copyCmd.SetArgs([]string{})
		copyCmd.SetOut(&stdOut)
		copyCmd.SetErr(&stdErr)

		path, err = tempFile("fixtures/vault.bin")
		Expect(err).ShouldNot(HaveOccurred())

		viper.Set("password.path", path)
		viper.Set("PASS_SESSION", "MRfbladUgDxLHvVWbxUjQUiZQykqiNcK")
	})

	Context("copy existing password to cpliboard", func() {
		It("copied password to clipboard", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for which you want to get your password:")
				c.SendLine("exist@mail.com")
				c.ExpectEOF()
			}()

			err = copyCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())
			//password, err := clipboard.ReadAll()
			//Expect(err).ShouldNot(HaveOccurred())
			//Expect(password).To(Equal("gMdLasZIGAEmDSCprqFkZQSAnjzeZzUP"))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("get non-existing password", func() {
		It("should throw password", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for which you want to get your password:")
				c.SendLine("nonexist@mail.com")
				c.ExpectEOF()
			}()

			err = copyCmd.Execute()
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
			err = copyCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})
