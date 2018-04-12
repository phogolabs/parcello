package parcel_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcel"
	"github.com/phogolabs/parcel/fake"
)

var _ = Describe("Generator", func() {
	var (
		generator *parcel.Generator
		bundle    *parcel.Bundle
	)

	BeforeEach(func() {
		bundle = &parcel.Bundle{
			Name: "bundle",
			Body: parcel.NewBufferCloser([]byte("hello")),
		}

		generator = &parcel.Generator{
			Config: &parcel.GeneratorConfig{},
		}
	})

	It("writes the bundle to the destination successfully", func() {
		buffer := &bytes.Buffer{}
		Expect(generator.WriteTo(buffer, bundle)).To(Succeed())
		Expect(buffer.String()).To(ContainSubstring("func init()"))
		Expect(buffer.String()).To(ContainSubstring("parcel.AddResource"))
		Expect(buffer.String()).NotTo(ContainSubstring("// Auto-generated"))
	})

	Context("when include API documentation is enabled", func() {
		BeforeEach(func() {
			generator.Config.InlcudeDocs = true
		})

		It("includes the documentation", func() {
			buffer := &bytes.Buffer{}
			Expect(generator.WriteTo(buffer, bundle)).To(Succeed())
			Expect(buffer.String()).To(ContainSubstring("func init()"))
			Expect(buffer.String()).To(ContainSubstring("parcel.AddResource"))
			Expect(buffer.String()).To(ContainSubstring("// Auto-generated"))
		})
	})

	Context("when writing the bundle fails", func() {
		It("returns the error", func() {
			buffer := &fake.ReadWriteCloser{}
			buffer.WriteReturns(0, fmt.Errorf("Oh no!"))
			Expect(generator.WriteTo(buffer, bundle)).To(MatchError("Oh no!"))
		})
	})
})
