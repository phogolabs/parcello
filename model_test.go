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

		It("creates a directory node", func() {
			dir := parcel.NewNodeDir("jack")
			Expect(dir.Name()).To(Equal("jack"))
			Expect(dir.IsDir()).To(BeTrue())
		})

		It("creates a file node", func() {
			file := parcel.NewNodeFile("jack", []byte{1})
			Expect(file.Name()).To(Equal("jack"))
			Expect(file.Size()).To(Equal(int64(1)))
		})

		It("returns the Name successfully", func() {
			Expect(node.Name()).To(BeEmpty())
		})

		It("returns the Size successfully", func() {
			file := parcel.NewNodeFile("jack", []byte{1})
			Expect(file.Name()).To(Equal("jack"))
			Expect(file.Size()).To(Equal(int64(1)))
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
		var buffer *parcel.Buffer

		Context("when the node is file", func() {
			BeforeEach(func() {
				node := parcel.NewNodeFile("sample.txt", []byte("hello"))
				buffer = parcel.NewBuffer(node)
			})

			It("reads successfully", func() {
				data, err := ioutil.ReadAll(buffer)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal("hello"))
			})

			It("writes successfully", func() {
				fmt.Fprintf(buffer, ",jack")

				data, err := ioutil.ReadAll(buffer)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal("hello,jack"))
			})

			It("closes successfully", func() {
				Expect(buffer.Close()).To(Succeed())
			})

			It("seeks successfully", func() {
				n, err := buffer.Seek(0, 0)
				Expect(err).To(BeNil())
				Expect(n).To(BeZero())
			})

			It("returns as string", func() {
				Expect(buffer.String()).To(Equal("hello"))
			})

			It("reads the directory fails", func() {
				files, err := buffer.Readdir(-1)
				Expect(err).To(MatchError("Not supported"))
				Expect(files).To(HaveLen(0))
			})

			It("returns the information successfully", func() {
				info, err := buffer.Stat()
				Expect(err).To(BeNil())
				Expect(info.IsDir()).To(BeFalse())
				Expect(info.Name()).To(Equal("sample.txt"))
			})
		})

		Context("when the node is directory", func() {
			BeforeEach(func() {
				child1 := parcel.NewNodeFile("sample.txt", []byte("hello"))
				child2 := parcel.NewNodeFile("report.txt", []byte("world"))
				node := parcel.NewNodeDir("documents", child1, child2)
				buffer = parcel.NewBuffer(node)
			})

			It("reads the directory successfully", func() {
				files, err := buffer.Readdir(-1)
				Expect(err).To(BeNil())
				Expect(files).To(HaveLen(2))

				info := files[0]
				Expect(info.Name()).To(Equal("sample.txt"))

				info = files[1]
				Expect(info.Name()).To(Equal("report.txt"))
			})

			Context("when the n is 1", func() {
				It("reads the directory successfully", func() {
					files, err := buffer.Readdir(1)
					Expect(err).To(BeNil())
					Expect(files).To(HaveLen(1))

					info := files[0]
					Expect(info.Name()).To(Equal("sample.txt"))
				})
			})

			It("returns the information successfully", func() {
				info, err := buffer.Stat()
				Expect(err).To(BeNil())
				Expect(info.IsDir()).To(BeTrue())
				Expect(info.Name()).To(Equal("documents"))
			})
		})
	})
})
