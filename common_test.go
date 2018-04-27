package parcello_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Common", func() {
	Describe("NewManager", func() {
		It("creates a new manager successfully", func() {
			manager := parcello.NewManager()
			Expect(manager).NotTo(BeNil())
			_, ok := manager.(*parcello.ResourceManager)
			Expect(ok).To(BeTrue())
		})

		Context("when dev mode is enabled", func() {
			BeforeEach(func() {
				os.Setenv("PARCELLO_DEV_ENABLED", "1")
			})

			AfterEach(func() {
				os.Unsetenv("PARCELLO_DEV_ENABLED")
			})

			It("creates a new dir manager", func() {
				manager := parcello.NewManager()
				Expect(manager).NotTo(BeNil())
				dir, ok := manager.(parcello.Dir)
				Expect(ok).To(BeTrue())
				Expect(string(dir)).To(Equal("."))
			})

			Context("when the directory is provided", func() {
				BeforeEach(func() {
					os.Setenv("PARCELLO_RESOURCE_DIR", "./root")
				})

				AfterEach(func() {
					os.Unsetenv("PRACELLO_RESOURCE_DIR")
				})

				It("creates a new dir manager", func() {
					manager := parcello.NewManager()
					Expect(manager).NotTo(BeNil())
					dir, ok := manager.(parcello.Dir)
					Expect(ok).To(BeTrue())
					Expect(string(dir)).To(Equal("./root"))
				})
			})
		})
	})
})
