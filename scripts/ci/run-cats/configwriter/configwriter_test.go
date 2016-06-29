package configwriter_test

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/configwriter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func randomBool() bool {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(2) == 0
}

var _ = Describe("Configwriter", func() {
	It("Generates a config object", func() {
		config := configwriter.GenerateConfigFromEnv()
		Expect(config).NotTo(BeNil())
	})

	Context("When valid CF envvars are set", func() {

		expectedApi := "api.example.com" + "_" + time.Now().String()
		expectedAdminUser := "admin_user" + "_" + time.Now().String()
		expectedPassword := "admin_password" + "_" + time.Now().String()
		expectedAppsDomain := "apps_domain" + "_" + time.Now().String()
		expectedSkipSslValidation := randomBool()
		expectedUseHttp := randomBool()
		expectedExistingUser := "existing_user" + "_" + time.Now().String()
		expectedExistingUserPassword := "expected_existing_user_password" + "_" + time.Now().String()

		BeforeEach(func() {
			os.Setenv("CF_API", expectedApi)
			os.Setenv("CF_ADMIN_USER", expectedAdminUser)
			os.Setenv("CF_ADMIN_PASSWORD", expectedPassword)
			os.Setenv("CF_APPS_DOMAIN", expectedAppsDomain)
			os.Setenv("SKIP_SSL_VALIDATION", strconv.FormatBool(expectedSkipSslValidation))
			os.Setenv("USE_HTTP", strconv.FormatBool(expectedUseHttp))
			os.Setenv("EXISTING_USER", expectedExistingUser)
			os.Setenv("EXISTING_USER_PASSWORD", expectedExistingUserPassword)
		})

		AfterEach(func() {
			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
			os.Unsetenv("SKIP_SSL_VALIDATION")
			os.Unsetenv("USE_HTTP")
			os.Unsetenv("EXISTING_USER")
			os.Unsetenv("EXISTING_USER_PASSWORD")
		})

		It("Generates a config object with the correct CF env variables set", func() {
			config := configwriter.GenerateConfigFromEnv()
			Expect(config).NotTo(BeNil())
			Expect(config.Api).To(Equal(expectedApi))
			Expect(config.AdminUser).To(Equal(expectedAdminUser))
			Expect(config.AdminPassword).To(Equal(expectedPassword))
			Expect(config.AppsDomain).To(Equal(expectedAppsDomain))
			Expect(config.SkipSslValidation).To(Equal(expectedSkipSslValidation))
			Expect(config.UseHttp).To(Equal(expectedUseHttp))
			Expect(config.ExistingUser).To(Equal(expectedExistingUser))
			Expect(config.ExistingUserPassword).To(Equal(expectedExistingUserPassword))
		})

		It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to true if 'ExistingUser' is provided", func() {
			expectedExistingUser := "existing_user" + "_" + time.Now().String()
			os.Setenv("EXISTING_USER", expectedExistingUser)
			config := configwriter.GenerateConfigFromEnv()
			Expect(config.UseExistingUser).To(Equal(true))
			Expect(config.KeepUserAtSuiteEnd).To(Equal(true))
		})

		It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to false if 'ExistingUser' is not provided", func() {
			os.Unsetenv("EXISTING_USER")
			config := configwriter.GenerateConfigFromEnv()
			Expect(config.UseExistingUser).To(Equal(false))
			Expect(config.KeepUserAtSuiteEnd).To(Equal(false))
		})
	})

	It("Uses the correct keynames when marshalling to json", func() {
		configJson, err := json.Marshal(configwriter.GenerateConfigFromEnv())
		Expect(err).NotTo(HaveOccurred())
		Expect(string(configJson)).To(MatchJSON(`{
																							"api": "",
																							"admin_user": "",
																							"admin_password": "",
																							"apps_domain": "",
																							"skip_ssl_validation": false,
																							"use_http": false,
																							"existing_user": "",
																							"use_existing_user": false,
																							"keep_user_at_suite_end": false,
																							"existing_user_password": ""
																							}`))
	})
})
