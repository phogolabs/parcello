package embedo_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/embedo"
)

var _ = Describe("Resource", func() {
	var (
		manager  *embedo.ResourceManager
		resource *embedo.Resource
	)

	BeforeEach(func() {
		resource = &embedo.Resource{}
	})

	JustBeforeEach(func() {
		manager = embedo.Open(resource)
	})

	Describe("Open", func() {
		Context("when the resource is empty", func() {
			It("returns an error", func() {
				file, err := manager.Open("migration.sql")
				Expect(file).To(BeNil())
				Expect(err).To(MatchError("File 'migration.sql' not found"))
			})
		})

		Context("when the resource has embedded content", func() {
			BeforeEach(func() {
				resource.Add("migration.sql", []byte("hello"))
				resource.Add("/root/etc/passwd", []byte("swordfish"))
			})

			It("returns the resource successfully", func() {
				file, err := manager.Open("migration.sql")
				Expect(file).NotTo(BeNil())
				Expect(err).NotTo(HaveOccurred())

				data, err := ioutil.ReadAll(file)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal("hello"))
			})

			Context("when the file is in subdirectory", func() {
				It("returns the resource successfully", func() {
					file, err := manager.Open("/root/etc/passwd")
					Expect(file).NotTo(BeNil())
					Expect(err).NotTo(HaveOccurred())

					data, err := ioutil.ReadAll(file)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(data)).To(Equal("swordfish"))
				})
			})

			Context("when is trying to open a directory", func() {
				It("returns an error", func() {
					file, err := manager.Open("/root/")
					Expect(file).To(BeNil())
					Expect(err).To(MatchError("Cannot open directory '/root/'"))
				})
			})

			Context("when the file with the requested name does not exist", func() {
				It("returns an error", func() {
					file, err := manager.Open("/home/root/migration.sql")
					Expect(file).To(BeNil())
					Expect(err).To(MatchError("File '/home/root/migration.sql' not found"))
				})
			})
		})
	})

	Describe("Walk", func() {
		Context("when the resource is empty", func() {
			It("returns an error", func() {
				err := manager.Walk("/documents", func(path string, info os.FileInfo, err error) error {
					return nil
				})

				Expect(err).To(MatchError("Directory '/documents' not found"))
			})
		})

		FContext("when the resource has hierarchy of directories and files", func() {
			BeforeEach(func() {
				resource.Add("/migration.sql", []byte("hello"))
				resource.Add("/root/etc/passwd", []byte("swordfish"))
			})

			It("walks through all of them", func() {
				paths := []string{}
				err := manager.Walk("/", func(path string, info os.FileInfo, err error) error {
					paths = append(paths, path)
					return nil
				})

				Expect(paths).To(HaveLen(5))
				Expect(paths[0]).To(Equal("/"))
				Expect(paths[1]).To(Equal("/migration.sql"))
				Expect(paths[2]).To(Equal("/root"))
				Expect(paths[3]).To(Equal("/root/etc"))
				Expect(paths[4]).To(Equal("/root/etc/passwd"))
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when the start node is file", func() {
				It("walks through the file only", func() {
					cnt := 0
					err := manager.Walk("/migration.sql", func(path string, info os.FileInfo, err error) error {
						cnt = cnt + 1
						Expect(path).To(Equal("/migration.sql"))
						Expect(info.Name()).To(Equal("migration.sql"))
						Expect(info.Size()).NotTo(BeZero())
						return nil
					})

					Expect(err).NotTo(HaveOccurred())
					Expect(cnt).To(Equal(1))
				})
			})

			It("walks through all of root children", func() {
				cnt := 0
				paths := []string{}
				err := manager.Walk("/root", func(path string, info os.FileInfo, err error) error {
					paths = append(paths, path)
					cnt = cnt + 1
					return nil
				})

				Expect(paths).To(HaveLen(3))
				Expect(paths[0]).To(Equal("/root"))
				Expect(paths[1]).To(Equal("/root/etc"))
				Expect(paths[2]).To(Equal("/root/etc/passwd"))
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when the walker returns an error", func() {
				It("returns the error", func() {
					err := manager.Walk("/root", func(path string, info os.FileInfo, err error) error {
						return fmt.Errorf("Oh no!")
					})

					Expect(err).To(MatchError("Oh no!"))
				})
			})
		})
	})
})
