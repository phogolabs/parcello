package parcel_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcel"
	"github.com/phogolabs/parcel/fake"
)

var _ = Describe("Generator", func() {
	var (
		generator  *parcel.Generator
		bundle     *parcel.Bundle
		buffer     *parcel.Buffer
		fileSystem *fake.FileSystem
	)

	BeforeEach(func() {
		buffer = parcel.NewBuffer()

		fileSystem = &fake.FileSystem{}
		fileSystem.OpenFileReturns(buffer, nil)

		bundle = &parcel.Bundle{
			Name: "bundle",
			Body: parcel.NewBufferWith([]byte("hello")),
		}

		generator = &parcel.Generator{
			FileSystem: fileSystem,
			Config: &parcel.GeneratorConfig{
				Package: "mypackage",
			},
		}
	})

	It("writes the bundle to the destination successfully", func() {
		Expect(generator.Compose(bundle)).To(Succeed())
		Expect(fileSystem.OpenFileCallCount()).To(Equal(1))

		filename, flag, mode := fileSystem.OpenFileArgsForCall(0)
		Expect(filename).To(Equal("bundle.go"))
		Expect(flag).To(Equal(os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
		Expect(mode).To(Equal(os.FileMode(0600)))

		Expect(buffer.String()).To(ContainSubstring("package mypackage"))
		Expect(buffer.String()).To(ContainSubstring("func init()"))
		Expect(buffer.String()).To(ContainSubstring("parcel.AddResource"))
		Expect(buffer.String()).NotTo(ContainSubstring("// Auto-generated"))
	})

	Context("when include API documentation is enabled", func() {
		BeforeEach(func() {
			generator.Config.InlcudeDocs = true
		})

		It("includes the documentation", func() {
			Expect(generator.Compose(bundle)).To(Succeed())
			Expect(buffer.String()).To(ContainSubstring("package mypackage"))
			Expect(buffer.String()).To(ContainSubstring("func init()"))
			Expect(buffer.String()).To(ContainSubstring("parcel.AddResource"))
			Expect(buffer.String()).To(ContainSubstring("// Auto-generated"))
		})
	})

	Context("when the file system fails", func() {
		It("returns the error", func() {
			fileSystem.OpenFileReturns(nil, fmt.Errorf("Oh no!"))
			Expect(generator.Compose(bundle)).To(MatchError("Oh no!"))
		})
	})

	Context("when writing the bundle fails", func() {
		It("returns the error", func() {
			buffer := &fake.File{}
			buffer.WriteReturns(0, fmt.Errorf("Oh no!"))
			fileSystem.OpenFileReturns(buffer, nil)

			Expect(generator.Compose(bundle)).To(MatchError("Oh no!"))
		})
	})
})
