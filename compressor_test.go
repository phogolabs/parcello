package parcello_test

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/parcello/fake"
)

var _ = Describe("ZipCompressor", func() {
	var (
		compressor *parcello.ZipCompressor
	)

	BeforeEach(func() {
		compressor = &parcello.ZipCompressor{
			Config: &parcello.CompressorConfig{
				Logger:   GinkgoWriter,
				Filename: "bundle",
				Recurive: true,
			},
		}
	})

	It("compresses a given hierarchy", func() {
		fileSystem := parcello.Dir("./fixture")

		bundle, err := compressor.Compress(fileSystem)
		Expect(err).To(BeNil())
		Expect(bundle).NotTo(BeNil())
		Expect(bundle.Name).To(Equal("bundle"))

		reader, err := zip.NewReader(strings.NewReader(string(bundle.Body)), int64(len(bundle.Body)))
		Expect(err).To(BeNil())

		Expect(reader.File).To(HaveLen(4))
		Expect(reader.File[0].Name).To(Equal("resource/reports/2018.txt"))
		Expect(reader.File[1].Name).To(Equal("resource/scripts/schema.sql"))
		Expect(reader.File[2].Name).To(Equal("resource/templates/html/index.html"))
		Expect(reader.File[3].Name).To(Equal("resource/templates/yml/schema.yml"))
	})

	Context("whene ingore pattern is provided", func() {
		It("ignores that files", func() {
			compressor.Config.IgnorePatterns = []string{"*/**/*.txt"}
			fileSystem := parcello.Dir("./fixture")

			bundle, err := compressor.Compress(fileSystem)
			Expect(err).To(BeNil())
			Expect(bundle).NotTo(BeNil())
			Expect(bundle.Name).To(Equal("bundle"))

			reader, err := zip.NewReader(strings.NewReader(string(bundle.Body)), int64(len(bundle.Body)))
			Expect(err).To(BeNil())

			Expect(reader.File).To(HaveLen(3))
			Expect(reader.File[0].Name).To(Equal("resource/scripts/schema.sql"))
			Expect(reader.File[1].Name).To(Equal("resource/templates/html/index.html"))
			Expect(reader.File[2].Name).To(Equal("resource/templates/yml/schema.yml"))
		})

		Context("when the pattern is directory", func() {
			It("ignores the directory and its files", func() {
				compressor.Config.IgnorePatterns = []string{"resource/templates/**/*"}
				fileSystem := parcello.Dir("./fixture")

				bundle, err := compressor.Compress(fileSystem)
				Expect(err).To(BeNil())
				Expect(bundle).NotTo(BeNil())
				Expect(bundle.Name).To(Equal("bundle"))

				reader, err := zip.NewReader(strings.NewReader(string(bundle.Body)), int64(len(bundle.Body)))
				Expect(err).To(BeNil())
				Expect(reader.File[0].Name).To(Equal("resource/reports/2018.txt"))
				Expect(reader.File[1].Name).To(Equal("resource/scripts/schema.sql"))
			})
		})
	})

	Context("when the pattern is invalid", func() {
		It("returns an error", func() {
			compressor.Config.IgnorePatterns = []string{"[*"}
			fileSystem := parcello.Dir("./fixture")

			bundle, err := compressor.Compress(fileSystem)
			Expect(err).To(MatchError("syntax error in pattern"))
			Expect(bundle).To(BeNil())
		})
	})

	Context("when the recursion is disabled", func() {
		It("does not go through the hierarchy", func() {
			compressor.Config.Recurive = false

			fileSystem := parcello.Dir("./fixture")

			bundle, err := compressor.Compress(fileSystem)
			Expect(err).To(BeNil())
			Expect(bundle).To(BeNil())
		})
	})

	Context("when opening file fails", func() {
		It("return the error", func() {
			fileSystem := &fake.FileSystem{}
			fileSystem.WalkStub = parcello.Dir("./fixture").Walk
			fileSystem.OpenFileReturns(nil, fmt.Errorf("Oh no!"))

			binary, err := compressor.Compress(fileSystem)
			Expect(err).To(MatchError("Oh no!"))
			Expect(binary).To(BeNil())
		})
	})

	Context("when the walker returns an nil file info", func() {
		It("return the error", func() {
			fileSystem := &fake.FileSystem{}
			fileSystem.WalkStub = func(dir string, fn filepath.WalkFunc) error {
				return fn("/", nil, nil)
			}

			binary, err := compressor.Compress(fileSystem)
			Expect(err).To(BeNil())
			Expect(binary).To(BeNil())
		})
	})

	Context("when the walker callback has an error", func() {
		It("return the error", func() {
			fileSystem := &fake.FileSystem{}
			fileSystem.WalkStub = func(dir string, fn filepath.WalkFunc) error {
				return fn("path", nil, fmt.Errorf("Oh no!"))
			}

			binary, err := compressor.Compress(fileSystem)
			Expect(err).To(MatchError("Oh no!"))
			Expect(binary).To(BeNil())
		})
	})

	Context("when the traversing fails", func() {
		It("return the error", func() {
			fileSystem := &fake.FileSystem{}
			fileSystem.WalkReturns(fmt.Errorf("Oh no!"))

			binary, err := compressor.Compress(fileSystem)
			Expect(err).To(MatchError("Oh no!"))
			Expect(binary).To(BeNil())
		})
	})
})
