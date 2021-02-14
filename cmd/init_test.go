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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("Init", func() {
	var (
		c       *expect.Console
		state   *vt10x.State
		err     error
		initCmd *cobra.Command
		stdOut  bytes.Buffer
		stdErr  bytes.Buffer

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
		cmd.Client = &mongo.Client{}
		initCmd = cmd.NewInitCmd()

		initCmd.SetArgs([]string{})
		initCmd.SetOut(&stdOut)
		initCmd.SetErr(&stdErr)
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
				c.SendLine("test-dummy@mail.com")
				c.ExpectString("Enter your password:")
				c.SendLine("test")
				c.ExpectString("password should be longer than 8 characters")
				c.SendLine("test-dummy")
				c.ExpectString("Enter again your password:")
				c.SendLine("test-dummy")
				c.ExpectEOF()
			}()

			mockError := errors.New("Mock error")
			mockDatabase.EXPECT().Find("dummytest-dummy@mail.com", gomock.Any()).Return(bson.M{}, mockError)
			mockDatabase.EXPECT().Insert("dummytest-dummy@mail.com", gomock.Any(), gomock.Any())

			err = initCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("when account is present in db", func() {
		It("failed to insert account in db", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your email address:")
				c.SendLine("dummy")
				c.ExpectString("email should be longer than 8 characters")
				c.SendLine("test-dummy@mail.com")
				c.ExpectString("Enter your password:")
				c.SendLine("test")
				c.ExpectString("password should be longer than 8 characters")
				c.SendLine("test-dummy")
				c.ExpectString("Enter again your password:")
				c.SendLine("test-dummy")
				c.ExpectEOF()
			}()

			expectedErr := errors.New("failed to insert data to db")
			mockError := errors.New("Mock error")
			mockDatabase.EXPECT().Find("dummytest-dummy@mail.com", gomock.Any()).Return(bson.M{}, mockError)
			mockDatabase.EXPECT().Insert("dummytest-dummy@mail.com", gomock.Any(), gomock.Any()).Return(expectedErr)

			err = initCmd.Execute()
			Expect(err).To(Equal(expectedErr))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

})
