package parcel_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcel"
	"github.com/phogolabs/parcel/fake"
)

var _ = Describe("Emitter", func() {
	var (
		emitter    *parcel.Emitter
		composer   *fake.Composer
		compressor *fake.Compressor
		fileSystem *fake.FileSystem
		resource   *parcel.BufferCloser
		bundle     *parcel.Bundle
	)

	BeforeEach(func() {
		resource = parcel.NewBufferCloser([]byte("data"))

		bundle = &parcel.Bundle{
			Name:   "resource",
			Body:   parcel.NewBufferCloser([]byte("resource")),
			Length: 20,
		}

		compressor = &fake.Compressor{}
		compressor.CompressReturns(bundle, nil)

		composer = &fake.Composer{}

		fileSystem = &fake.FileSystem{}
		fileSystem.OpenFileReturns(resource, nil)

		emitter = &parcel.Emitter{
			Logger:     GinkgoWriter,
			Compressor: compressor,
			Composer:   composer,
			FileSystem: fileSystem,
		}
	})

	It("emits the provided source successfully", func() {
		Expect(emitter.Emit()).To(Succeed())
		Expect(compressor.CompressCallCount()).To(Equal(1))
		Expect(compressor.CompressArgsForCall(0)).To(Equal(fileSystem))
		Expect(fileSystem.OpenFileCallCount()).To(Equal(1))

		filename, mode, perm := fileSystem.OpenFileArgsForCall(0)
		Expect(filename).To(Equal("resource.go"))
		Expect(int(mode)).To(Equal(int(os.O_WRONLY | os.O_CREATE | os.O_TRUNC)))
		Expect(int(perm)).To(Equal(0600))

		Expect(composer.WriteToCallCount()).To(Equal(1))
		r, b := composer.WriteToArgsForCall(0)
		Expect(r).To(Equal(resource))
		Expect(b).To(Equal(bundle))
	})

	Context("when the bundle length is zero", func() {
		It("does not compose it", func() {
			bundle.Length = 0
			Expect(emitter.Emit()).To(Succeed())
			Expect(compressor.CompressCallCount()).To(Equal(1))
			Expect(compressor.CompressArgsForCall(0)).To(Equal(fileSystem))
			Expect(fileSystem.OpenFileCallCount()).To(Equal(0))
			Expect(composer.WriteToCallCount()).To(Equal(0))
		})
	})

	Context("when the compressor fails", func() {
		It("returns the error", func() {
			compressor.CompressReturns(nil, fmt.Errorf("Oh no!"))
			Expect(emitter.Emit()).To(MatchError("Oh no!"))
		})
	})

	Context("when the file system fails", func() {
		It("returns the error", func() {
			fileSystem.OpenFileReturns(nil, fmt.Errorf("Oh no!"))
			Expect(emitter.Emit()).To(MatchError("Oh no!"))
		})
	})

	Context("when the composer fails", func() {
		It("returns the error", func() {
			composer.WriteToReturns(fmt.Errorf("Oh no!"))
			Expect(emitter.Emit()).To(MatchError("Oh no!"))
		})
	})
})
