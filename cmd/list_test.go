package cmd_test

import (
	"bytes"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/nnachevv/passmag/cmd"
	"github.com/nnachevv/passmag/crypt"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ = Describe("List", func() {
	var (
		c       *expect.Console
		state   *vt10x.State
		err     error
		path    string
		listCmd *cobra.Command
		stdOut  bytes.Buffer
		stdErr  bytes.Buffer
	)

	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		cmd.Stdio = terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
		cmd.Crypt = crypt.Crypt{}

		listCmd = cmd.NewListCmd()
		listCmd.SetArgs([]string{})
		listCmd.SetOut(&stdOut)
		listCmd.SetErr(&stdErr)

		path, err = tempFile("fixtures/vault.bin")
		Expect(err).ShouldNot(HaveOccurred())

		viper.Set("password.path", path)
		viper.Set("PASS_SESSION", "MRfbladUgDxLHvVWbxUjQUiZQykqiNcK")
	})

	Context("list all password", func() {
		It("should print all passwords", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectEOF()
			}()

			err = listCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})
