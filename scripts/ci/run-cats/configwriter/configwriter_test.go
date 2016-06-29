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
			Expect(config).NotTo(BeNil())
			Expect(config.UseExistingUser).To(Equal(false))
			Expect(config.KeepUserAtSuiteEnd).To(Equal(false))
		})

		Context("When Valid Persistent App envvars are provided", func() {
			expectedPersistedAppHost := "PERSISTENT_APP_HOST" + "_" + time.Now().String()
			expectedPersistedAppSpace := "PERSISTENT_APP_SPACE" + "_" + time.Now().String()
			expectedPersistedAppOrg := "PERSISTENT_APP_ORG" + "_" + time.Now().String()
			expectedPersistedAppQuotaName := "PERSISTENT_APP_QUOTA_NAME" + "_" + time.Now().String()

			BeforeEach(func() {
				os.Setenv("PERSISTENT_APP_HOST", expectedPersistedAppHost)
				os.Setenv("PERSISTENT_APP_SPACE", expectedPersistedAppSpace)
				os.Setenv("PERSISTENT_APP_ORG", expectedPersistedAppOrg)
				os.Setenv("PERSISTENT_APP_QUOTA_NAME", expectedPersistedAppQuotaName)
			})
			AfterEach(func() {
				os.Unsetenv("PERSISTENT_APP_HOST")
				os.Unsetenv("PERSISTENT_APP_SPACE")
				os.Unsetenv("PERSISTENT_APP_ORG")
				os.Unsetenv("PERSISTENT_APP_QUOTA_NAME")
			})

			It("Generates a config object with the correct persistent variables set", func() {
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.PersistentAppHost).To(Equal(expectedPersistedAppHost))
				Expect(config.PersistentAppSpace).To(Equal(expectedPersistedAppSpace))
				Expect(config.PersistentAppOrg).To(Equal(expectedPersistedAppOrg))
				Expect(config.PersistentAppQuotaName).To(Equal(expectedPersistedAppQuotaName))
			})

			It("Does not generate a config object with PersistentAppHost if it is not set", func() {
				os.Unsetenv("PERSISTENT_APP_HOST")

				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.PersistentAppSpace).To(Equal(expectedPersistedAppSpace))
				Expect(config.PersistentAppOrg).To(Equal(expectedPersistedAppOrg))
				Expect(config.PersistentAppQuotaName).To(Equal(expectedPersistedAppQuotaName))

				Expect(config.PersistentAppHost).To(Equal(""))
			})

			It("Does not generate a config object with PersistentAppSpace if it is not set", func() {
				os.Unsetenv("PERSISTENT_APP_SPACE")

				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.PersistentAppHost).To(Equal(expectedPersistedAppHost))
				Expect(config.PersistentAppOrg).To(Equal(expectedPersistedAppOrg))
				Expect(config.PersistentAppQuotaName).To(Equal(expectedPersistedAppQuotaName))

				Expect(config.PersistentAppSpace).To(Equal(""))
			})

			It("Does not generate a config object with PersistentAppOrg if it is not set", func() {
				os.Unsetenv("PERSISTENT_APP_ORG")

				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.PersistentAppHost).To(Equal(expectedPersistedAppHost))
				Expect(config.PersistentAppSpace).To(Equal(expectedPersistedAppSpace))
				Expect(config.PersistentAppQuotaName).To(Equal(expectedPersistedAppQuotaName))

				Expect(config.PersistentAppOrg).To(Equal(""))
			})

			It("Does not generate a config object with PersistentAppQuotaName if it is not set", func() {
				os.Unsetenv("PERSISTENT_APP_QUOTA_NAME")

				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.PersistentAppHost).To(Equal(expectedPersistedAppHost))
				Expect(config.PersistentAppSpace).To(Equal(expectedPersistedAppSpace))
				Expect(config.PersistentAppOrg).To(Equal(expectedPersistedAppOrg))

				Expect(config.PersistentAppQuotaName).To(Equal(""))
			})

		})

		Context("When one or more of the allowed buildpacks is provided", func() {
			expectedStaticBuildpackName := "STATICFILE_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedJavaBuildpackName := "JAVA_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedRubyBuildpackName := "Ruby_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedNodeJsBuildpackName := "NODEJS_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedGoBuildpackName := "GO_BUILDPACK_NAME" + "_" + time.Now().String()

			expectedPythonBuildpackName := "PYTHON_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedPhpBuildpackName := "PHP_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedBinaryBuildpackName := "BINARY_BUILDPACK_NAME" + "_" + time.Now().String()

			BeforeEach(func() {
				os.Setenv("STATICFILE_BUILDPACK_NAME", expectedStaticBuildpackName)
				os.Setenv("JAVA_BUILDPACK_NAME", expectedJavaBuildpackName)
				os.Setenv("RUBY_BUILDPACK_NAME", expectedRubyBuildpackName)
				os.Setenv("NODEJS_BUILDPACK_NAME", expectedNodeJsBuildpackName)
				os.Setenv("GO_BUILDPACK_NAME", expectedGoBuildpackName)
				os.Setenv("PYTHON_BUILDPACK_NAME", expectedPythonBuildpackName)
				os.Setenv("PHP_BUILDPACK_NAME", expectedPhpBuildpackName)
				os.Setenv("BINARY_BUILDPACK_NAME", expectedBinaryBuildpackName)
			})
			AfterEach(func() {
				os.Unsetenv("STATICFILE_BUILDPACK_NAME")
				os.Unsetenv("JAVA_BUILDPACK_NAME")
				os.Unsetenv("RUBY_BUILDPACK_NAME")
				os.Unsetenv("NODEJS_BUILDPACK_NAME")
				os.Unsetenv("GO_BUILDPACK_NAME")
				os.Unsetenv("PYTHON_BUILDPACK_NAME")
				os.Unsetenv("PHP_BUILDPACK_NAME")
				os.Unsetenv("BINARY_BUILDPACK_NAME")
			})

			It("Generates a config object with the correct buildpack variables set", func() {
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.StaticBuildpackName).To(Equal(expectedStaticBuildpackName))
				Expect(config.JavaBuildpackName).To(Equal(expectedJavaBuildpackName))
				Expect(config.RubyBuildpackName).To(Equal(expectedRubyBuildpackName))
				Expect(config.NodeJsBuildpackName).To(Equal(expectedNodeJsBuildpackName))
				Expect(config.GoBuildpackName).To(Equal(expectedGoBuildpackName))
				Expect(config.PythonBuildpackName).To(Equal(expectedPythonBuildpackName))
				Expect(config.PhpBuildpackName).To(Equal(expectedPhpBuildpackName))
				Expect(config.BinaryBuildpackName).To(Equal(expectedBinaryBuildpackName))
			})

			It("Does not generate a config object with StaticBuildpack if it is not set", func() {
				os.Unsetenv("STATICFILE_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.StaticBuildpackName).To(Equal(""))
			})
			It("Does not generate a config object with JavaBuildpack if it is not set", func() {
				os.Unsetenv("JAVA_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.JavaBuildpackName).To(Equal(""))
			})
			It("Does not generate a config object with RubyBuildpack if it is not set", func() {
				os.Unsetenv("RUBY_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.RubyBuildpackName).To(Equal(""))
			})
			It("Does not generate a config object with NodeJsBuildpack if it is not set", func() {
				os.Unsetenv("NODEJS_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.NodeJsBuildpackName).To(Equal(""))
			})
			It("Does not generate a config object with GoBuildpack if it is not set", func() {
				os.Unsetenv("GO_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.GoBuildpackName).To(Equal(""))
			})
			It("Does not generate a config object with PythonBuildpack if it is not set", func() {
				os.Unsetenv("PYTHON_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.PythonBuildpackName).To(Equal(""))
			})
			It("Does not generate a config object with PHPBuildpack if it is not set", func() {
				os.Unsetenv("PHP_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.PhpBuildpackName).To(Equal(""))
			})
			It("Does not generate a config object with BinaryBuildpack if it is not set", func() {
				os.Unsetenv("BINARY_BUILDPACK_NAME")
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.BinaryBuildpackName).To(Equal(""))
			})
		})
	})

	Context("When no env variables are provided", func() {
		It("Only renders the required keynames when marshalling to json", func() {
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
																							"existing_user_password": "",
																							"staticfile_buildpack_name": "",
																							"java_buildpack_name": "",
																							"ruby_buildpack_name": "",
																							"nodejs_buildpack_name": "",
																							"go_buildpack_name": "",
																							"python_buildpack_name": "",
																							"php_buildpack_name": "",
																							"binary_buildpack_name": "",
																							"persistent_app_host": "",
																							"persistent_app_space": "",
																							"persistent_app_org": "",
																							"persistent_app_quota_name": ""
																							}`))
		})
	})

	Context("When all the env variable are provided", func() {
		It("Renders a valid integration_config when marshalling to json", func() {

		})
	})
})
