package parcello_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

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
		It("opens the root successfully", func() {
			file, err := manager.Open("/")
			Expect(file).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		Context("when the resource is empty", func() {
			It("returns an error", func() {
				file, err := manager.Open("/migration.sql")
				Expect(file).To(BeNil())
				Expect(err).To(MatchError("Directory does not exist"))
			})
		})

		Context("when the file is directory", func() {
			It("returns an error", func() {
				file, err := manager.Open("/resource/reports")
				Expect(file).NotTo(BeNil())
				Expect(err).To(BeNil())
			})
		})

		Context("when the global resource is empty", func() {
			It("returns an error", func() {
				file, err := parcello.Open("migration.sql")
				Expect(file).To(BeNil())
				Expect(err).To(MatchError("Directory does not exist"))
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

		Context("when the file is open more than once for read", func() {
			It("does not change the mod time", func() {
				file, err := manager.Open("/resource/reports/2018.txt")
				Expect(file).NotTo(BeNil())
				Expect(err).NotTo(HaveOccurred())

				info, err := file.Stat()
				Expect(err).NotTo(HaveOccurred())

				file, err = manager.Open("/resource/reports/2018.txt")
				Expect(file).NotTo(BeNil())
				Expect(err).NotTo(HaveOccurred())

				info2, err := file.Stat()
				Expect(err).NotTo(HaveOccurred())

				Expect(info.ModTime()).To(Equal(info2.ModTime()))
			})
		})

		It("returns a readonly resource", func() {
			file, err := manager.Open("/resource/reports/2018.txt")
			Expect(file).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			_, err = fmt.Fprintln(file.(io.Writer), "hello")
			Expect(err).To(MatchError("File is read-only"))
		})

		Context("when the file with the requested name does not exist", func() {
			It("returns an error", func() {
				file, err := manager.Open("/resource/migration.sql")
				Expect(file).To(BeNil())
				Expect(err).To(MatchError("open /resource/migration.sql: file does not exist"))
			})
		})
	})

	Describe("OpenFile", func() {
		Context("when the file does not exist", func() {
			It("creates the file", func() {
				file, err := manager.OpenFile("/resource/secrets.txt", os.O_CREATE, 0600)
				Expect(file).NotTo(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the file is directory", func() {
			It("returns an error", func() {
				file, err := manager.OpenFile("/resource/reports", os.O_RDWR, 0600)
				Expect(file).To(BeNil())
				Expect(err).To(MatchError("open /resource/reports: Is directory"))
			})
		})

		Context("when the file exists", func() {
			It("truncs the file content", func() {
				file, err := manager.OpenFile("/resource/reports/2018.txt", os.O_CREATE|os.O_TRUNC, 0600)
				Expect(file).NotTo(BeNil())
				Expect(err).NotTo(HaveOccurred())

				data, err := ioutil.ReadAll(file)
				Expect(err).NotTo(HaveOccurred())
				Expect(data).To(BeEmpty())
			})

			Context("when the file is open more than once for write", func() {
				It("does not change the mod time", func() {
					start := time.Now()

					file, err := manager.OpenFile("/resource/reports/2018.txt", os.O_WRONLY, 0600)
					Expect(file).NotTo(BeNil())
					Expect(err).NotTo(HaveOccurred())

					info, err := file.Stat()
					Expect(err).NotTo(HaveOccurred())
					modTime := info.ModTime()

					Expect(modTime.After(start)).To(BeTrue())
				})
			})

			Context("when the os.O_TRUNC flag is not provided", func() {
				It("returns an error", func() {
					file, err := manager.OpenFile("/resource/reports/2018.txt", os.O_CREATE, 0600)
					Expect(file).To(BeNil())
					Expect(err).To(MatchError("open /resource/reports/2018.txt: file already exists"))
				})
			})

			Context("when the file is open for append", func() {
				It("appends content successfully", func() {
					file, err := manager.OpenFile("/resource/reports/2018.txt", os.O_RDWR|os.O_APPEND, 0600)
					Expect(file).NotTo(BeNil())
					Expect(err).NotTo(HaveOccurred())

					_, err = fmt.Fprint(file, "hello")
					Expect(err).NotTo(HaveOccurred())

					_, err = file.Seek(0, os.SEEK_SET)
					Expect(err).NotTo(HaveOccurred())

					data, err := ioutil.ReadAll(file)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(data)).To(Equal("Report 2018\nhello"))
				})
			})

			Context("when the file is open for WRITE only", func() {
				Context("when we try to read", func() {
					It("returns an error", func() {
						file, err := manager.OpenFile("/resource/reports/2018.txt", os.O_WRONLY, 0600)
						Expect(file).NotTo(BeNil())
						Expect(err).NotTo(HaveOccurred())

						_, err = ioutil.ReadAll(file)
						Expect(err).To(MatchError("File is write-only"))
					})
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
