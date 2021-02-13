package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/golang/mock/gomock"
	"github.com/hinshun/vt10x"
	"github.com/nnachevv/passmag/cmd"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/mocks"
	"github.com/nnachevv/passmag/storage"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

var _ = Describe("Add", func() {
	var (
		c           *expect.Console
		state       *vt10x.State
		err         error
		path        string
		addCmd      *cobra.Command
		stdOut      bytes.Buffer
		stdErr      bytes.Buffer
		vaultPwd    []byte
		mockCtrl    *gomock.Controller
		mockMongoDB *mocks.MockMongoDatabase
	)

	BeforeEach(func() {
		c, state, err = vt10x.NewVT10XConsole()
		Expect(err).ShouldNot(HaveOccurred())
		cmd.Stdio = terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
		cmd.Crypt = crypt.Crypt{}
		mockCtrl = gomock.NewController(GinkgoT())
		mockMongoDB = mocks.NewMockMongoDatabase(mockCtrl)
		cmd.MongoDB = mockMongoDB

		addCmd = cmd.NewAddCmd()
		addCmd.SetArgs([]string{})
		addCmd.SetOut(&stdOut)
		addCmd.SetErr(&stdErr)

		vaultPwd = argon2.IDKey([]byte("test-dummy"), []byte("MRfbladUgDxLHvVWbxUjQUiZQykqiNcK"), 1, 64*1024, 4, 32)
		path, err = tempFile("fixtures/vault.bin")
		Expect(err).ShouldNot(HaveOccurred())

		viper.Set("password.path", path)
		viper.Set("PASS_SESSION", "MRfbladUgDxLHvVWbxUjQUiZQykqiNcK")
	})
	AfterEach(func() {

		mockCtrl.Finish()
	})
	Context("When user set PASS_SESSION ,decline generation, pass right arguments for his passwords", func() {
		It("should add to vault his password", func() {
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for your password:")
				c.SendLine("dummy-name")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
				c.ExpectEOF()
			}()
			mockMongoDB.EXPECT().Insert("exist@mail.com", gomock.Any())
			err = addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())
			password, ok := getAddedPassword(path, "dummy-name", vaultPwd)
			Expect(ok).Should(BeTrue())
			Expect(string(password)).To(Equal("dummy-password"))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("want password to be generated, pass name for his password", func() {
		It("should add to vault his password", func() {
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for your password:")
				c.SendLine("dummy-random")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("y")
				c.ExpectEOF()

			}()
			mockMongoDB.EXPECT().Insert("exist@mail.com", gomock.Any())
			err = addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())
			_, ok := getAddedPassword(path, "dummy-random", vaultPwd)
			Expect(ok).Should(BeTrue())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("trying to add already exist password and decline editing already existing password", func() {
		It("should add to vault his password", func() {

			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for your password:")
				c.SendLine("exist")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
				c.ExpectString("This name with password already exist! Do you want to edit name with newly password")
				c.SendLine("N")
				c.ExpectEOF()

			}()
			mockMongoDB.EXPECT().Insert("exist@mail.com", gomock.Any())
			err = addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())
			password, ok := getAddedPassword(path, "exist", vaultPwd)
			Expect(ok).Should(BeTrue())
			Expect(string(password)).To(Equal("dummy-password"))

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("trying to add already exist password and edit", func() {
		It("should add to vault his password", func() {
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your master password:")
				c.SendLine("test-dummy")
				c.ExpectString("Enter name for your password:")
				c.SendLine("exist@mail.com")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
				c.ExpectString("This name with password already exist! Do you want to edit name with newly password")
				c.SendLine("y")
				c.ExpectEOF()
			}()
			mockMongoDB.EXPECT().Insert("exist@mail.com", gomock.Any())
			err = addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())
			password, ok := getAddedPassword(path, "exist@mail.com", vaultPwd)
			Expect(ok).Should(BeTrue())
			Expect(string(password)).To(Equal("dummy-password"))

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
			err = addCmd.Execute()
			Expect(err).Should(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})

// creates a new temporary file
func tempFile(path string) (string, error) {
	tar, err := os.Open(path)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(tar)
	Expect(err).ShouldNot(HaveOccurred())
	defer tar.Close()

	file, err := ioutil.TempFile(os.TempDir(), "fixture-file")
	Expect(err).ShouldNot(HaveOccurred())

	_, err = file.Write(bytes)
	Expect(err).ShouldNot(HaveOccurred())

	err = file.Sync()
	Expect(err).ShouldNot(HaveOccurred())

	_, err = file.Seek(0, io.SeekStart)
	Expect(err).ShouldNot(HaveOccurred())

	return file.Name(), nil
}

func getAddedPassword(path string, name string, vaultPwd []byte) (string, bool) {
	vaultData, err := cmd.Crypt.DecryptFile(path, vaultPwd)
	Expect(err).ShouldNot(HaveOccurred())

	var s storage.Storage

	err = json.Unmarshal(vaultData, &s)
	Expect(err).ShouldNot(HaveOccurred())
	password, ok := s.Passwords[name]

	return string(password), ok
}
