package main_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	Context("when no envvars are set", func() {

		It("prints an error regarding missing 'required' env vars", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say(`Missing required environment variables:
CF_API
CF_ADMIN_USER
CF_ADMIN_PASSWORD
CF_APPS_DOMAIN
EXISTING_USER
EXISTING_USER_PASSWORD`,
			))
		})
	})

	Context("when some envvars are set", func() {
		BeforeEach(func() {
			os.Setenv("CF_API", "cf_api_value")
		})
		AfterEach(func() {
			os.Unsetenv("CF_API")
		})

		It("prints only the missing 'required' env vars", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say(`Missing required environment variables:
CF_ADMIN_USER
CF_ADMIN_PASSWORD
CF_APPS_DOMAIN
EXISTING_USER
EXISTING_USER_PASSWORD`,
			))

		})
	})
})
