package embedo_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/embedo"
)

var _ = Describe("Generator", func() {
	var (
		generator *embedo.Generator
		dir       string
	)

	BeforeEach(func() {
		var err error

		dir, err = ioutil.TempDir("", "gom")
		Expect(err).To(BeNil())

		generator = &embedo.Generator{
			FileSystem: embedo.Dir(dir),
			Config: &embedo.GeneratorConfig{
				Recurive:    true,
				InlcudeDocs: false,
			},
		}

		Expect(os.MkdirAll(filepath.Join(dir, "sub"), 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(dir, "sub", "more.sql"), []byte("more"), 0600)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(dir, "script.sql"), []byte("hello"), 0600)).To(Succeed())
	})

	It("generates the embedded resources successfully", func() {
		Expect(generator.Generate("resource")).To(Succeed())

		path := filepath.Join(dir, "resource.go")
		Expect(path).To(BeARegularFile())

		data, err := ioutil.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(ContainSubstring("more.sql"))
		Expect(string(data)).To(ContainSubstring("script.sql"))
		Expect(string(data)).NotTo(ContainSubstring("// Auto-generated"))
	})

	Context("when the recursion is disabled", func() {
		BeforeEach(func() {
			generator.Config.Recurive = false
		})

		It("generates the embedded resources on root level successfully", func() {
			Expect(generator.Generate("resource")).To(Succeed())

			path := filepath.Join(dir, "resource.go")
			Expect(path).To(BeARegularFile())

			data, err := ioutil.ReadFile(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).NotTo(ContainSubstring("more.sql"))
			Expect(string(data)).To(ContainSubstring("script.sql"))
			Expect(string(data)).NotTo(ContainSubstring("// Auto-generated"))
		})
	})

	Context("when the documentation should be included", func() {
		BeforeEach(func() {
			generator.Config.InlcudeDocs = true
		})

		It("generates the embedded resources successfully", func() {
			Expect(generator.Generate("resource")).To(Succeed())

			path := filepath.Join(dir, "resource.go")
			Expect(path).To(BeARegularFile())

			data, err := ioutil.ReadFile(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring("script.sql"))
			Expect(string(data)).To(ContainSubstring("// Auto-generated"))
		})
	})

	Context("when the file system fails", func() {
		It("returns an error", func() {
			generator.FileSystem = embedo.Dir("/database")
			Expect(generator.Generate("resource")).To(MatchError("Directory does not exist"))
		})
	})
})
