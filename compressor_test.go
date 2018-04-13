package parcel_test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcel"
	"github.com/phogolabs/parcel/fake"
)

var _ = Describe("TarGZipCompressor", func() {
	var (
		compressor *parcel.TarGZipCompressor
	)

	BeforeEach(func() {
		compressor = &parcel.TarGZipCompressor{
			Config: &parcel.CompressorConfig{
				Logger:   GinkgoWriter,
				Name:     "bundle",
				Recurive: true,
			},
		}
	})

	It("compresses a given hierarchy", func() {
		fileSystem := parcel.Dir("./fixture")

		bundle, err := compressor.Compress(fileSystem)
		Expect(err).To(BeNil())
		Expect(bundle).NotTo(BeNil())
		Expect(bundle.Name).To(Equal("bundle"))

		gzipper, err := gzip.NewReader(bundle.Body)
		Expect(err).To(BeNil())

		reader := tar.NewReader(gzipper)

		header, err := reader.Next()
		Expect(err).To(BeNil())
		Expect(header.Name).To(Equal("resource/reports/2018.txt"))

		header, err = reader.Next()
		Expect(err).To(BeNil())
		Expect(header.Name).To(Equal("resource/scripts/schema.sql"))

		header, err = reader.Next()
		Expect(err).To(BeNil())
		Expect(header.Name).To(Equal("resource/templates/html/index.html"))

		header, err = reader.Next()
		Expect(err).To(BeNil())
		Expect(header.Name).To(Equal("resource/templates/yml/schema.yml"))

		header, err = reader.Next()
		Expect(header).To(BeNil())
		Expect(err).To(MatchError("unexpected EOF"))

		Expect(bundle.Body.Close()).To(Succeed())
	})

	Context("whene ingore pattern is provided", func() {
		It("ignores that files", func() {
			compressor.Config.IgnorePatterns = []string{"*/**/*.txt"}
			fileSystem := parcel.Dir("./fixture")

			bundle, err := compressor.Compress(fileSystem)
			Expect(err).To(BeNil())
			Expect(bundle).NotTo(BeNil())
			Expect(bundle.Name).To(Equal("bundle"))

			gzipper, err := gzip.NewReader(bundle.Body)
			Expect(err).To(BeNil())

			reader := tar.NewReader(gzipper)

			header, err := reader.Next()
			Expect(err).To(BeNil())
			Expect(header.Name).To(Equal("resource/scripts/schema.sql"))

			header, err = reader.Next()
			Expect(err).To(BeNil())
			Expect(header.Name).To(Equal("resource/templates/html/index.html"))

			header, err = reader.Next()
			Expect(err).To(BeNil())
			Expect(header.Name).To(Equal("resource/templates/yml/schema.yml"))

			header, err = reader.Next()
			Expect(header).To(BeNil())
			Expect(err).To(MatchError("unexpected EOF"))

			Expect(bundle.Body.Close()).To(Succeed())
		})

		Context("when the pattern is directory", func() {
			It("ignores the directory and its files", func() {
				compressor.Config.IgnorePatterns = []string{"resource/templates/**/*"}
				fileSystem := parcel.Dir("./fixture")

				bundle, err := compressor.Compress(fileSystem)
				Expect(err).To(BeNil())
				Expect(bundle).NotTo(BeNil())
				Expect(bundle.Name).To(Equal("bundle"))

				gzipper, err := gzip.NewReader(bundle.Body)
				Expect(err).To(BeNil())

				reader := tar.NewReader(gzipper)

				header, err := reader.Next()
				Expect(err).To(BeNil())
				Expect(header.Name).To(Equal("resource/reports/2018.txt"))

				header, err = reader.Next()
				Expect(err).To(BeNil())
				Expect(header.Name).To(Equal("resource/scripts/schema.sql"))

				header, err = reader.Next()
				Expect(header).To(BeNil())
				Expect(err).To(MatchError("unexpected EOF"))

				Expect(bundle.Body.Close()).To(Succeed())
			})
		})
	})

	Context("when the pattern is invalid", func() {
		It("returns an error", func() {
			compressor.Config.IgnorePatterns = []string{"[*"}
			fileSystem := parcel.Dir("./fixture")

			bundle, err := compressor.Compress(fileSystem)
			Expect(err).To(MatchError("syntax error in pattern"))
			Expect(bundle).To(BeNil())
		})
	})

	Context("when the recursion is disabled", func() {
		It("does not go through the hierarchy", func() {
			compressor.Config.Recurive = false

			fileSystem := parcel.Dir("./fixture")

			bundle, err := compressor.Compress(fileSystem)
			Expect(err).To(BeNil())
			Expect(bundle).NotTo(BeNil())
			Expect(bundle.Name).To(Equal("bundle"))

			data, err := ioutil.ReadAll(bundle.Body)
			Expect(bundle.Body.Close()).To(Succeed())
			Expect(data).To(HaveLen(10))
		})
	})

	Context("when opening file fails", func() {
		It("return the error", func() {
			fileSystem := &fake.FileSystem{}
			fileSystem.WalkStub = parcel.Dir("./fixture").Walk
			fileSystem.OpenFileReturns(nil, fmt.Errorf("Oh no!"))

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
