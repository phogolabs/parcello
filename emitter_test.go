package parcello_test

import (
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/parcello/fake"
)

var _ = Describe("Emitter", func() {
	var (
		emitter    *parcello.Emitter
		composer   *fake.Composer
		compressor *fake.Compressor
		fileSystem *fake.FileSystem
		resource   *parcello.ResourceFile
		bundle     *parcello.Bundle
	)

	BeforeEach(func() {
		data := []byte("data")
		node := &parcello.Node{
			Name:    "resource",
			Content: &data,
			Mutex:   &sync.RWMutex{},
		}

		resource = parcello.NewResourceFile(node)

		bundle = &parcello.Bundle{
			Info: &parcello.BundleInfo{
				Name:  "resource",
				Count: 20,
			},
			Body: []byte("resource"),
		}

		compressor = &fake.Compressor{}
		compressor.CompressStub = func(ctx *parcello.CompressorContext) (*parcello.BundleInfo, error) {
			ctx.Writer.Write(bundle.Body)
			return bundle.Info, nil
		}

		composer = &fake.Composer{}

		fileSystem = &fake.FileSystem{}
		fileSystem.OpenFileReturns(resource, nil)

		emitter = &parcello.Emitter{
			Logger:     GinkgoWriter,
			Compressor: compressor,
			Composer:   composer,
			FileSystem: fileSystem,
		}
	})

	It("emits the provided source successfully", func() {
		Expect(emitter.Emit()).To(Succeed())
		Expect(compressor.CompressCallCount()).To(Equal(1))

		ctx := compressor.CompressArgsForCall(0)
		Expect(ctx.FileSystem).To(Equal(fileSystem))

		Expect(composer.ComposeCallCount()).To(Equal(1))
		Expect(composer.ComposeArgsForCall(0)).To(Equal(bundle))
	})

	Context("when the bundle is nil", func() {
		It("does not compose it", func() {
			compressor.CompressReturns(nil, nil)

			Expect(emitter.Emit()).To(Succeed())
			Expect(compressor.CompressCallCount()).To(Equal(1))
			ctx := compressor.CompressArgsForCall(0)
			Expect(ctx.FileSystem).To(Equal(fileSystem))
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
