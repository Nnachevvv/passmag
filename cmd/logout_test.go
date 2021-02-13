package cmd_test

import (
	"bytes"

	"github.com/nnachevv/passmag/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("List", func() {
	var (
		err       error
		logoutCmd *cobra.Command
		stdOut    bytes.Buffer
		stdErr    bytes.Buffer
	)

	BeforeEach(func() {
		logoutCmd = cmd.NewLogoutCmd()
		logoutCmd.SetArgs([]string{})
		logoutCmd.SetOut(&stdOut)
		logoutCmd.SetErr(&stdErr)
	})

	Context("try to delete vault", func() {
		It("give error - vault do not exist", func() {
			err = logoutCmd.Execute()
			Expect(err).Should(HaveOccurred())

		})
	})
})
