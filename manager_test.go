package parcello_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Manager", func() {
	var (
		manager  *parcello.Manager
		resource []byte
	)

	BeforeEach(func() {
		manager = &parcello.Manager{}

		var err error
		compressor := parcello.TarGZipCompressor{
			Config: &parcello.CompressorConfig{
				Logger:   ioutil.Discard,
				Filename: "bundle",
				Recurive: true,
			},
		}

		bundle, err := compressor.Compress(parcello.Dir("./fixture"))
		Expect(err).NotTo(HaveOccurred())

		resource = bundle.Body
	})

	JustBeforeEach(func() {
		Expect(manager.Add(resource)).To(Succeed())
	})

	Describe("Add", func() {
		Context("when the resource is added second time", func() {
			It("returns an error", func() {
				Expect(manager.Add(resource)).To(MatchError("Invalid path: 'resource/reports/2018.txt'"))
			})
		})

		Context("when the resource is not tar.gz", func() {
			It("returns an error", func() {
				Expect(manager.Add([]byte("lol"))).To(MatchError("unexpected EOF"))
			})

			It("panics", func() {
				Expect(func() { parcello.AddResource([]byte("lol")) }).To(Panic())
			})
		})
	})

	Describe("Root", func() {
		It("returns a valid sub-manager", func() {
			group, err := manager.Root("/resource")
			Expect(err).To(BeNil())

			file, err := group.Open("/reports/2018.txt")
			Expect(file).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			data, err := ioutil.ReadAll(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal("Report 2018\n"))
		})

		Context("when group is a file not a directory", func() {
			It("returns an error", func() {
				group, err := manager.Root("/resource/reports/2018.txt")
				Expect(group).To(BeNil())
				Expect(err).To(MatchError("Resource hierarchy not found"))
			})
		})

		Context("when the manager is global", func() {
			It("returns a valid sub-manager", func() {
				parcello.AddResource(resource)
				group := parcello.Root("/resource")

				file, err := group.Open("/reports/2018.txt")
				Expect(file).NotTo(BeNil())
				Expect(err).NotTo(HaveOccurred())

				data, err := ioutil.ReadAll(file)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal("Report 2018\n"))
			})

			Context("when group is a file not a directory", func() {
				It("panics", func() {
					Expect(func() { parcello.Root("/resource/reports/2018.txt") }).To(Panic())
				})
			})
		})
	})

	Describe("Open", func() {
		Context("when the resource is empty", func() {
			It("returns an error", func() {
				file, err := manager.Open("migration.sql")
				Expect(file).To(BeNil())
				Expect(err).To(MatchError("File 'migration.sql' not found"))
			})
		})

		Context("when the global resource is empty", func() {
			It("returns an error", func() {
				file, err := parcello.Open("migration.sql")
				Expect(file).To(BeNil())
				Expect(err).To(MatchError("File 'migration.sql' not found"))
			})
		})

		It("returns the resource successfully", func() {
			file, err := manager.Open("/resource/reports/2018.txt")
			Expect(file).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			data, err := ioutil.ReadAll(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal("Report 2018\n"))
		})

		Context("when is trying to open a directory", func() {
			It("returns an error", func() {
				file, err := manager.Open("/resource/reports/")
				Expect(file).NotTo(BeNil())
				Expect(err).To(BeNil())
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

	Describe("Walk", func() {
		Context("when the resource is empty", func() {
			It("returns an error", func() {
				err := manager.Walk("/documents", func(path string, info os.FileInfo, err error) error {
					return nil
				})

				Expect(err).To(MatchError("Directory '/documents' not found"))
			})
		})

		Context("when the resource has hierarchy of directories and files", func() {
			It("walks through all of them", func() {
				paths := []string{}
				err := manager.Walk("/", func(path string, info os.FileInfo, err error) error {
					paths = append(paths, path)
					return nil
				})

				Expect(paths).To(HaveLen(11))
				Expect(paths[0]).To(Equal("/"))
				Expect(paths[1]).To(Equal("/resource"))
				Expect(paths[2]).To(Equal("/resource/reports"))
				Expect(paths[3]).To(Equal("/resource/reports/2018.txt"))
				Expect(paths[4]).To(Equal("/resource/scripts"))
				Expect(paths[5]).To(Equal("/resource/scripts/schema.sql"))
				Expect(paths[6]).To(Equal("/resource/templates"))
				Expect(paths[7]).To(Equal("/resource/templates/html"))
				Expect(paths[8]).To(Equal("/resource/templates/html/index.html"))
				Expect(paths[9]).To(Equal("/resource/templates/yml"))
				Expect(paths[10]).To(Equal("/resource/templates/yml/schema.yml"))
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when the start node is file", func() {
				It("walks through the file only", func() {
					cnt := 0
					err := manager.Walk("/resource/reports/2018.txt", func(path string, info os.FileInfo, err error) error {
						cnt = cnt + 1
						Expect(path).To(Equal("/resource/reports/2018.txt"))
						Expect(info.Name()).To(Equal("2018.txt"))
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
				err := manager.Walk("/resource/templates", func(path string, info os.FileInfo, err error) error {
					paths = append(paths, path)
					cnt = cnt + 1
					return nil
				})

				Expect(paths).To(HaveLen(5))
				Expect(paths[0]).To(Equal("/resource/templates"))
				Expect(paths[1]).To(Equal("/resource/templates/html"))
				Expect(paths[2]).To(Equal("/resource/templates/html/index.html"))
				Expect(paths[3]).To(Equal("/resource/templates/yml"))
				Expect(paths[4]).To(Equal("/resource/templates/yml/schema.yml"))
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when the walker returns an error", func() {
				It("returns the error", func() {
					err := manager.Walk("/resource", func(path string, info os.FileInfo, err error) error {
						return fmt.Errorf("Oh no!")
					})

					Expect(err).To(MatchError("Oh no!"))
				})

				Context("when the walk returns an error for sub-directory", func() {
					It("returns the error", func() {
						err := manager.Walk("/resource", func(path string, info os.FileInfo, err error) error {
							if path == "/resource/templates" {
								return fmt.Errorf("Oh no!")
							}
							return nil
						})

						Expect(err).To(MatchError("Oh no!"))
					})
				})
			})
		})
	})
})
