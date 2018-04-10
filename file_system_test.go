package embedo_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/embedo"
)

var _ = Describe("FileSystem", func() {
	var dir embedo.Dir

	BeforeEach(func() {
		path, err := ioutil.TempDir("", "gom_generator")
		Expect(err).To(BeNil())

		dir = embedo.Dir(path)
		Expect(ioutil.WriteFile(filepath.Join(path, "sample.txt"), []byte("test"), 0600)).To(Succeed())
	})

	Context("OpenFile", func() {
		It("opens a file successfully", func() {
			file, err := dir.OpenFile("sample.txt", os.O_RDONLY, 0)
			Expect(err).To(BeNil())

			content, err := ioutil.ReadAll(file)
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal("test"))
			Expect(file.Close()).To(Succeed())
		})

		Context("when the file does not exists", func() {
			It("returns an error", func() {
				file, err := dir.OpenFile("report.txt", os.O_RDONLY, 0)
				Expect(file).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})

	Context("Walk", func() {
		It("walks through the hierarchy successfully", func() {
			count := 0
			err := dir.Walk("/", func(path string, info os.FileInfo, err error) error {
				count = count + 1

				if info.IsDir() {
					Expect(path).To(Equal("."))
				} else {
					Expect(path).To(Equal("sample.txt"))
				}

				return nil
			})

			Expect(count).To(Equal(2))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the walking fails", func() {
			It("returns an error", func() {
				err := dir.Walk("/wrong", func(path string, info os.FileInfo, err error) error {
					return fmt.Errorf("Oh no!")
				})

				Expect(err).To(MatchError("Oh no!"))
			})
		})
	})
})
