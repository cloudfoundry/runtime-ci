package configwriter_test

import (
	"encoding/json"
	"os"
	"time"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/configwriter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configwriter", func() {
	It("Generates a config object", func() {
		config := configwriter.GenerateConfig("", "")
		Expect(config).NotTo(BeNil())
	})

	Context("When a valid CF_API envvar is set", func() {
		BeforeEach(func() {
			os.Setenv("CF_API", "api.example.com")
		})

		AfterEach(func() {
			os.Unsetenv("CF_API")
		})

		It("Generates a config object with the correct 'api'", func() {
			config := configwriter.GenerateConfigFromEnv()
			Expect(config).NotTo(BeNil())
			Expect(config.Api).To(Equal("api.example.com"))
		})

	})

	Context("When a valid CF_ADMIN_USER envvar is set", func() {
		expectedAdminUser := "admin_user" + "_" + time.Now().String()
		BeforeEach(func() {
			os.Setenv("CF_ADMIN_USER", expectedAdminUser)
		})

		AfterEach(func() {
			os.Unsetenv("CF_ADMIN_USER")
		})

		It("Generates a config object with the correct 'admin_user'", func() {
			config := configwriter.GenerateConfigFromEnv()
			Expect(config).NotTo(BeNil())
			Expect(config.AdminUser).To(Equal(expectedAdminUser))
		})
	})

	Context("When a valid CF admin_user property is passed", func() {

		It("Generates a config object with 'admin_user' set correctly", func() {
			expectedAdminUser := "admin_user" + "_" + time.Now().String()
			config := configwriter.GenerateConfig("", expectedAdminUser)
			Expect(config).NotTo(BeNil())
			Expect(config.AdminUser).To(Equal(expectedAdminUser))
		})
	})
	It("Uses the correct keynames when marshalling to json", func() {
		configJson, err := json.Marshal(configwriter.GenerateConfigFromEnv())
		Expect(err).NotTo(HaveOccurred())
		Expect(string(configJson)).To(Equal(
			`{
"api": "",
"admin_user": ""
}`,
		))
	})
})
