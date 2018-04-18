package parcel_test

import (
	"fmt"

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
		resource   *parcel.Buffer
		bundle     *parcel.Bundle
	)

	BeforeEach(func() {
		resource = parcel.NewBuffer(parcel.NewNodeFile("resource", []byte("data")))

		bundle = &parcel.Bundle{
			Name:   "resource",
			Body:   []byte("resource"),
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

		Expect(composer.ComposeCallCount()).To(Equal(1))
		Expect(composer.ComposeArgsForCall(0)).To(Equal(bundle))
	})

	Context("when the bundle is nil", func() {
		It("does not compose it", func() {
			compressor.CompressReturns(nil, nil)

			Expect(emitter.Emit()).To(Succeed())
			Expect(compressor.CompressCallCount()).To(Equal(1))
			Expect(compressor.CompressArgsForCall(0)).To(Equal(fileSystem))
			Expect(composer.ComposeCallCount()).To(BeZero())
		})
	})

	Context("when the compressor fails", func() {
		It("returns the error", func() {
			compressor.CompressReturns(nil, fmt.Errorf("Oh no!"))
			Expect(emitter.Emit()).To(MatchError("Oh no!"))
		})
	})

	Context("when the composer fails", func() {
		It("returns the error", func() {
			composer.ComposeReturns(fmt.Errorf("Oh no!"))
			Expect(emitter.Emit()).To(MatchError("Oh no!"))
		})
	})
})
