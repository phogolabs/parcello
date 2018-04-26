package parcello_test

import (
	"fmt"
	"os"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/parcello/fake"
)

var _ = Describe("Bundler", func() {
	var (
		bundler    *parcello.Bundler
		compressor *fake.Compressor
		source     *fake.FileSystem
		target     *fake.FileSystem
		binary     *fake.File
		binaryInfo *parcello.ResourceFileInfo
		ctx        *parcello.BundlerContext
	)

	BeforeEach(func() {
		content := []byte("file")
		binaryInfo = &parcello.ResourceFileInfo{
			Node: &parcello.Node{
				Mutex:   &sync.RWMutex{},
				IsDir:   false,
				Content: &content,
			},
		}

		binary = &fake.File{}
		binary.StatReturns(binaryInfo, nil)
		compressor = &fake.Compressor{}
		source = &fake.FileSystem{}
		target = &fake.FileSystem{}
		target.OpenFileReturns(binary, nil)
		bundler = &parcello.Bundler{
			Logger:     GinkgoWriter,
			Compressor: compressor,
			FileSystem: source,
		}

		ctx = &parcello.BundlerContext{
			Name:       "app",
			FileSystem: target,
		}
	})

	It("bunles the binary successfully", func() {
		Expect(bundler.Bundle(ctx)).To(Succeed())
		Expect(target.OpenFileCallCount()).To(Equal(1))

		name, opts, perm := target.OpenFileArgsForCall(0)
		Expect(name).To(Equal(ctx.Name))
		Expect(opts).To(Equal(os.O_WRONLY))
		Expect(perm).To(Equal(os.FileMode(0600)))

		Expect(compressor.CompressCallCount()).To(Equal(1))

		cctx := compressor.CompressArgsForCall(0)
		Expect(cctx.FileSystem).To(Equal(source))
		Expect(cctx.Writer).To(Equal(binary))
		Expect(cctx.Offset).To(Equal(binaryInfo.Size()))
	})

	Context("when opening the binary fails", func() {
		BeforeEach(func() {
			target.OpenFileReturns(nil, fmt.Errorf("Oh no!"))
		})

		It("returns an error", func() {
			Expect(bundler.Bundle(ctx)).To(MatchError("Oh no!"))
		})
	})

	Context("when getting the binary information fails", func() {
		BeforeEach(func() {
			binary.StatReturns(nil, fmt.Errorf("Oh no!"))
		})

		It("returns an error", func() {
			Expect(bundler.Bundle(ctx)).To(MatchError("Oh no!"))
		})
	})

	Context("when the target is directory", func() {
		BeforeEach(func() {
			binaryInfo.Node.IsDir = true
		})

		It("returns an error", func() {
			Expect(bundler.Bundle(ctx)).To(MatchError("'app' is not a regular file"))
		})
	})

	Context("when the compressor fails", func() {
		BeforeEach(func() {
			compressor.CompressReturns(nil, fmt.Errorf("Oh no!"))
		})

		It("returns an error", func() {
			Expect(bundler.Bundle(ctx)).To(MatchError("Oh no!"))
		})
	})
})
