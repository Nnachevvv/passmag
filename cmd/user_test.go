package cmd_test

import (
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
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

var _ = Describe("User", func() {
	var (
		c         *expect.Console
		state     *vt10x.State
		err       error
		mockCtrl  *gomock.Controller
		mockCrypt *mocks.MockCrypter
	)
	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		cmd.Stdio = terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
		mockCtrl = gomock.NewController(GinkgoT())
		mockCrypt = mocks.NewMockCrypter(mockCtrl)
		cmd.Crypt = mockCrypt
		// set pass_session because pass_session has been set in previous tests
		viper.Set("PASS_SESSION", "")

	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("with valid account", func() {
		It("contains account in db", func() {
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Please enter your session key:")
				c.SendLine("test-sessionkey")
				c.ExpectString("Enter your master password:")
				c.SendLine("master-password")
				c.ExpectEOF()
			}()

			vaultPwd := argon2.IDKey([]byte("master-password"), []byte("test-sessionkey"), 1, 64*1024, 4, 32)
			mockCrypt.EXPECT().DecryptFile(gomock.Any(), vaultPwd)
			user, err := cmd.EnterSession()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(user.VaultPwd).To(Equal(vaultPwd))
			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})
