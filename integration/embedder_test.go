package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Embedder", func() {
	var (
		cmd  *exec.Cmd
		dir  string
		args []string
		res  string
	)

	BeforeEach(func() {
		args = []string{}
	})

	JustBeforeEach(func() {
		var err error

		dir, err = ioutil.TempDir("", "gom")
		Expect(err).To(BeNil())

		cmd = exec.Command(embedoPath, append(args, "-d", "./database", "-pkg", "resource")...)
		cmd.Dir = dir

		path := filepath.Join(cmd.Dir, "/database")
		Expect(os.MkdirAll(path, 0700)).To(Succeed())

		path = filepath.Join(path, "main.sql")
		Expect(ioutil.WriteFile(path, []byte("main"), 0700)).To(Succeed())

		path = filepath.Join(cmd.Dir, "/database/command")
		Expect(os.MkdirAll(path, 0700)).To(Succeed())

		path = filepath.Join(path, "commands.sql")
		Expect(ioutil.WriteFile(path, []byte("command"), 0700)).To(Succeed())

		res = filepath.Join(cmd.Dir, "/database/resource.go")
	})

	It("generates resource on root level", func() {
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		Expect(session.Out).To(gbytes.Say("Embedding 'main.sql'"))
		Expect(session.Out).NotTo(gbytes.Say("Embedding 'command/commands.sql'"))
		Expect(res).To(BeARegularFile())

		data, err := ioutil.ReadFile(res)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(ContainSubstring("main.sql"))
		Expect(string(data)).To(ContainSubstring("//"))
		Expect(string(data)).NotTo(ContainSubstring("commands.sql"))
	})

	Context("when the documentation is disabled", func() {
		BeforeEach(func() {
			args = append(args, "-include-docs=false")
		})

		It("does not include API documentation", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).To(gbytes.Say("Embedding 'main.sql'"))
			Expect(res).To(BeARegularFile())

			data, err := ioutil.ReadFile(res)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring("main.sql"))
			Expect(string(data)).NotTo(ContainSubstring("//"))
			Expect(string(data)).NotTo(ContainSubstring("commands.sql"))
		})
	})

	Context("when quite model is enabled", func() {
		BeforeEach(func() {
			args = append(args, "-q")
		})

		It("does not print anything on stdout", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).NotTo(gbytes.Say("Embedding 'main.sql'"))
			Expect(session.Out).NotTo(gbytes.Say("Embedding 'command/commands.sql'"))
			Expect(res).To(BeARegularFile())

			data, err := ioutil.ReadFile(res)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring("main.sql"))
			Expect(string(data)).NotTo(ContainSubstring("commands.sql"))
		})
	})

	Context("when the recursion is enabled", func() {
		BeforeEach(func() {
			args = append(args, "-r")
		})

		It("generates resource for all directories", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).To(gbytes.Say("Embedding 'command/commands.sql'"))
			Expect(session.Out).To(gbytes.Say("Embedding 'main.sql'"))
			Expect(res).To(BeARegularFile())

			data, err := ioutil.ReadFile(res)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring("main.sql"))
			Expect(string(data)).To(ContainSubstring("commands.sql"))
		})
	})

	Context("when the directory is not provided", func() {
		It("returns an error", func() {
			cmd.Args = cmd.Args[:1]

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(101))

			Expect(session.Err).To(gbytes.Say("Directory is not provided"))
			Expect(res).NotTo(BeARegularFile())
		})
	})

	Context("when the package is not provided", func() {
		It("returns an error", func() {
			cmd = exec.Command(embedoPath, append(args, "-d", "./database")...)

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(101))

			Expect(session.Err).To(gbytes.Say("Package name is not provided"))
			Expect(res).NotTo(BeARegularFile())
		})
	})
})
