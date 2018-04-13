package parcel_test

import (
	"fmt"
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcel"
)

var _ = Describe("Model", func() {
	Describe("Node", func() {
		var node *parcel.Node

		BeforeEach(func() {
			node = &parcel.Node{}
		})

		It("returns the Name successfully", func() {
			Expect(node.Name()).To(BeEmpty())
		})

		It("returns the Size successfully", func() {
			Expect(node.Size()).To(BeZero())
		})

		It("returns the Mode successfully", func() {
			Expect(node.Mode()).To(BeZero())
		})

		It("returns the ModTime successfully", func() {
			Expect(node.ModTime()).To(BeTemporally("~", time.Now()))
		})

		It("returns the IsDir successfully", func() {
			Expect(node.IsDir()).To(BeFalse())
		})

		It("returns the Sys successfully", func() {
			Expect(node.Sys()).To(BeNil())
		})
	})

	Describe("Buffer", func() {
		It("reads successfully", func() {
			buffer := parcel.NewBuffer([]byte("hello"))
			data, err := ioutil.ReadAll(buffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal("hello"))
		})

		It("writes successfully", func() {
			buffer := parcel.NewBuffer([]byte("hello"))
			fmt.Fprintf(buffer, ",jack")

			data, err := ioutil.ReadAll(buffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal("hello,jack"))
		})

		It("closes successfully", func() {
			buffer := parcel.NewBuffer([]byte("hello"))
			Expect(buffer.Close()).To(Succeed())
		})

		It("seeks successfully", func() {
			buffer := parcel.NewBuffer([]byte("hello"))
			n, err := buffer.Seek(0, 0)
			Expect(err).To(BeNil())
			Expect(n).To(BeZero())
		})

		It("returns as string", func() {
			buffer := parcel.NewBuffer([]byte("hello"))
			Expect(buffer.String()).To(Equal("hello"))
		})
	})
})
