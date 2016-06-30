package configwriter_test

import (
	"encoding/json"
	"io/ioutil"
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
	Context("when env vars are not set", func() {
		It("returns an empty config object", func() {
			config := configwriter.GenerateConfigFromEnv()
			Expect(config).NotTo(BeNil())
			Expect(config.Api).To(Equal(""))
			Expect(config.AdminUser).To(Equal(""))
			Expect(config.AdminPassword).To(Equal(""))
			Expect(config.AppsDomain).To(Equal(""))
			Expect(config.SkipSslValidation).To(BeFalse())
			Expect(config.UseHttp).To(BeFalse())
			Expect(config.ExistingUser).To(Equal(""))
			Expect(config.ExistingUserPassword).To(Equal(""))
			Expect(config.Backend).To(Equal(""))
			Expect(config.PersistentAppHost).To(Equal(""))
			Expect(config.PersistentAppSpace).To(Equal(""))
			Expect(config.PersistentAppOrg).To(Equal(""))
			Expect(config.PersistentAppQuotaName).To(Equal(""))
			Expect(config.StaticBuildpackName).To(Equal(""))
			Expect(config.JavaBuildpackName).To(Equal(""))
			Expect(config.RubyBuildpackName).To(Equal(""))
			Expect(config.NodeJsBuildpackName).To(Equal(""))
			Expect(config.GoBuildpackName).To(Equal(""))
			Expect(config.PythonBuildpackName).To(Equal(""))
			Expect(config.PhpBuildpackName).To(Equal(""))
			Expect(config.BinaryBuildpackName).To(Equal(""))
		})
	})

	Context("When valid CF envvars are set", func() {
		var (
			expectedApi                   string
			expectedAdminUser             string
			expectedPassword              string
			expectedAppsDomain            string
			expectedSkipSslValidation     bool
			expectedUseHttp               bool
			expectedExistingUser          string
			expectedExistingUserPassword  string
			expectedBackend               string
			expectedPersistedAppHost      string
			expectedPersistedAppSpace     string
			expectedPersistedAppOrg       string
			expectedPersistedAppQuotaName string
			expectedDefaultTimeout        int
			expectedCfPushTimeout         int
			expectedLongCurlTimeout       int
			expectedBrokerStartTimeout    int
			expectedStaticBuildpackName   string
			expectedJavaBuildpackName     string
			expectedRubyBuildpackName     string
			expectedNodeJsBuildpackName   string
			expectedGoBuildpackName       string
			expectedPythonBuildpackName   string
			expectedPhpBuildpackName      string
			expectedBinaryBuildpackName   string
		)

		BeforeEach(func() {
			expectedApi = "api.example.com" + "_" + time.Now().String()
			expectedAdminUser = "admin_user" + "_" + time.Now().String()
			expectedPassword = "admin_password" + "_" + time.Now().String()
			expectedAppsDomain = "apps_domain" + "_" + time.Now().String()
			expectedSkipSslValidation = randomBool()
			expectedUseHttp = randomBool()
			expectedExistingUser = "existing_user" + "_" + time.Now().String()
			expectedExistingUserPassword = "expected_existing_user_password" + "_" + time.Now().String()
			expectedBackend = "expected_backend" + "_" + time.Now().String()
			expectedPersistedAppHost = "PERSISTENT_APP_HOST" + "_" + time.Now().String()
			expectedPersistedAppSpace = "PERSISTENT_APP_SPACE" + "_" + time.Now().String()
			expectedPersistedAppOrg = "PERSISTENT_APP_ORG" + "_" + time.Now().String()
			expectedPersistedAppQuotaName = "PERSISTENT_APP_QUOTA_NAME" + "_" + time.Now().String()

			rand.Seed(time.Now().UTC().UnixNano())
			expectedDefaultTimeout = rand.Intn(100)
			expectedCfPushTimeout = rand.Intn(100)
			expectedLongCurlTimeout = rand.Intn(100)
			expectedBrokerStartTimeout = rand.Intn(100)

			expectedStaticBuildpackName = "STATICFILE_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedJavaBuildpackName = "JAVA_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedRubyBuildpackName = "Ruby_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedNodeJsBuildpackName = "NODEJS_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedGoBuildpackName = "GO_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedPythonBuildpackName = "PYTHON_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedPhpBuildpackName = "PHP_BUILDPACK_NAME" + "_" + time.Now().String()
			expectedBinaryBuildpackName = "BINARY_BUILDPACK_NAME" + "_" + time.Now().String()

			os.Setenv("CF_API", expectedApi)
			os.Setenv("CF_ADMIN_USER", expectedAdminUser)
			os.Setenv("CF_ADMIN_PASSWORD", expectedPassword)
			os.Setenv("CF_APPS_DOMAIN", expectedAppsDomain)
			os.Setenv("SKIP_SSL_VALIDATION", strconv.FormatBool(expectedSkipSslValidation))
			os.Setenv("USE_HTTP", strconv.FormatBool(expectedUseHttp))
			os.Setenv("EXISTING_USER", expectedExistingUser)
			os.Setenv("EXISTING_USER_PASSWORD", expectedExistingUserPassword)
			os.Setenv("BACKEND", expectedBackend)
			os.Setenv("PERSISTENT_APP_HOST", expectedPersistedAppHost)
			os.Setenv("PERSISTENT_APP_SPACE", expectedPersistedAppSpace)
			os.Setenv("PERSISTENT_APP_ORG", expectedPersistedAppOrg)
			os.Setenv("PERSISTENT_APP_QUOTA_NAME", expectedPersistedAppQuotaName)
			os.Setenv("DEFAULT_TIMEOUT", strconv.Itoa(expectedDefaultTimeout))
			os.Setenv("CF_PUSH_TIMEOUT", strconv.Itoa(expectedCfPushTimeout))
			os.Setenv("LONG_CURL_TIMEOUT", strconv.Itoa(expectedLongCurlTimeout))
			os.Setenv("BROKER_START_TIMEOUT", strconv.Itoa(expectedBrokerStartTimeout))
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
			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
			os.Unsetenv("SKIP_SSL_VALIDATION")
			os.Unsetenv("USE_HTTP")
			os.Unsetenv("EXISTING_USER")
			os.Unsetenv("EXISTING_USER_PASSWORD")
			os.Unsetenv("BACKEND")
			os.Unsetenv("PERSISTENT_APP_HOST")
			os.Unsetenv("PERSISTENT_APP_SPACE")
			os.Unsetenv("PERSISTENT_APP_ORG")
			os.Unsetenv("PERSISTENT_APP_QUOTA_NAME")
			os.Unsetenv("DEFAULT_TIMEOUT")
			os.Unsetenv("CF_PUSH_TIMEOUT")
			os.Unsetenv("LONG_CURL_TIMEOUT")
			os.Unsetenv("BROKER_START_TIMEOUT")
			os.Unsetenv("STATICFILE_BUILDPACK_NAME")
			os.Unsetenv("JAVA_BUILDPACK_NAME")
			os.Unsetenv("RUBY_BUILDPACK_NAME")
			os.Unsetenv("NODEJS_BUILDPACK_NAME")
			os.Unsetenv("GO_BUILDPACK_NAME")
			os.Unsetenv("PYTHON_BUILDPACK_NAME")
			os.Unsetenv("PHP_BUILDPACK_NAME")
			os.Unsetenv("BINARY_BUILDPACK_NAME")
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
			Expect(config.Backend).To(Equal(expectedBackend))
			Expect(config.PersistentAppHost).To(Equal(expectedPersistedAppHost))
			Expect(config.PersistentAppSpace).To(Equal(expectedPersistedAppSpace))
			Expect(config.PersistentAppOrg).To(Equal(expectedPersistedAppOrg))
			Expect(config.PersistentAppQuotaName).To(Equal(expectedPersistedAppQuotaName))
			Expect(config.DefaultTimeout).To(Equal(expectedDefaultTimeout))
			Expect(config.CfPushTimeout).To(Equal(expectedCfPushTimeout))
			Expect(config.LongCurlTimeout).To(Equal(expectedLongCurlTimeout))
			Expect(config.BrokerStartTimeout).To(Equal(expectedBrokerStartTimeout))
			Expect(config.StaticBuildpackName).To(Equal(expectedStaticBuildpackName))
			Expect(config.JavaBuildpackName).To(Equal(expectedJavaBuildpackName))
			Expect(config.RubyBuildpackName).To(Equal(expectedRubyBuildpackName))
			Expect(config.NodeJsBuildpackName).To(Equal(expectedNodeJsBuildpackName))
			Expect(config.GoBuildpackName).To(Equal(expectedGoBuildpackName))
			Expect(config.PythonBuildpackName).To(Equal(expectedPythonBuildpackName))
			Expect(config.PhpBuildpackName).To(Equal(expectedPhpBuildpackName))
			Expect(config.BinaryBuildpackName).To(Equal(expectedBinaryBuildpackName))
		})

		Context("when 'ExistingUser' is provided", func() {
			BeforeEach(func() {
				os.Setenv("EXISTING_USER", expectedExistingUser)
			})

			It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to true if 'ExistingUser' is provided", func() {
				config := configwriter.GenerateConfigFromEnv()
				Expect(config.UseExistingUser).To(Equal(true))
				Expect(config.KeepUserAtSuiteEnd).To(Equal(true))
			})
		})

		Context("when 'ExistingUser' is not provided", func() {
			BeforeEach(func() {
				os.Unsetenv("EXISTING_USER")
			})

			It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to false if 'ExistingUser' is not provided", func() {
				config := configwriter.GenerateConfigFromEnv()
				Expect(config).NotTo(BeNil())
				Expect(config.UseExistingUser).To(Equal(false))
				Expect(config.KeepUserAtSuiteEnd).To(Equal(false))
			})
		})
	})

	Describe("marshaling the struct", func() {
		It("does not render optional keys if their values are empty", func() {
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

		Context("when any env variables are provided", func() {
			BeforeEach(func() {
				os.Setenv("CF_API", "non-empty-value")
				os.Setenv("CF_ADMIN_USER", "non-empty-value")
				os.Setenv("CF_ADMIN_PASSWORD", "non-empty-value")
				os.Setenv("CF_APPS_DOMAIN", "non-empty-value")
				os.Setenv("SKIP_SSL_VALIDATION", "true")
				os.Setenv("USE_HTTP", "true")
				os.Setenv("EXISTING_USER", "non-empty-value")
				os.Setenv("EXISTING_USER_PASSWORD", "non-empty-value")
				os.Setenv("BACKEND", "non-empty-value")
				os.Setenv("PERSISTENT_APP_HOST", "non-empty-value")
				os.Setenv("PERSISTENT_APP_SPACE", "non-empty-value")
				os.Setenv("PERSISTENT_APP_ORG", "non-empty-value")
				os.Setenv("PERSISTENT_APP_QUOTA_NAME", "non-empty-value")
				os.Setenv("DEFAULT_TIMEOUT", "1")
				os.Setenv("CF_PUSH_TIMEOUT", "1")
				os.Setenv("LONG_CURL_TIMEOUT", "1")
				os.Setenv("BROKER_START_TIMEOUT", "1")
				os.Setenv("STATICFILE_BUILDPACK_NAME", "non-empty-value")
				os.Setenv("JAVA_BUILDPACK_NAME", "non-empty-value")
				os.Setenv("RUBY_BUILDPACK_NAME", "non-empty-value")
				os.Setenv("NODEJS_BUILDPACK_NAME", "non-empty-value")
				os.Setenv("GO_BUILDPACK_NAME", "non-empty-value")
				os.Setenv("PYTHON_BUILDPACK_NAME", "non-empty-value")
				os.Setenv("PHP_BUILDPACK_NAME", "non-empty-value")
				os.Setenv("BINARY_BUILDPACK_NAME", "non-empty-value")
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
				os.Unsetenv("BACKEND")
				os.Unsetenv("PERSISTENT_APP_HOST")
				os.Unsetenv("PERSISTENT_APP_SPACE")
				os.Unsetenv("PERSISTENT_APP_ORG")
				os.Unsetenv("PERSISTENT_APP_QUOTA_NAME")
				os.Unsetenv("DEFAULT_TIMEOUT")
				os.Unsetenv("CF_PUSH_TIMEOUT")
				os.Unsetenv("LONG_CURL_TIMEOUT")
				os.Unsetenv("BROKER_START_TIMEOUT")
				os.Unsetenv("STATICFILE_BUILDPACK_NAME")
				os.Unsetenv("JAVA_BUILDPACK_NAME")
				os.Unsetenv("RUBY_BUILDPACK_NAME")
				os.Unsetenv("NODEJS_BUILDPACK_NAME")
				os.Unsetenv("GO_BUILDPACK_NAME")
				os.Unsetenv("PYTHON_BUILDPACK_NAME")
				os.Unsetenv("PHP_BUILDPACK_NAME")
				os.Unsetenv("BINARY_BUILDPACK_NAME")
			})

			It("renders the variables in the integration_config", func() {
				configJson, err := json.Marshal(configwriter.GenerateConfigFromEnv())
				Expect(err).NotTo(HaveOccurred())
				Expect(string(configJson)).To(MatchJSON(`{
																							"api": "non-empty-value",
																							"admin_user": "non-empty-value",
																							"admin_password": "non-empty-value",
																							"apps_domain": "non-empty-value",
																							"skip_ssl_validation": true,
																							"use_http": true,
																							"existing_user": "non-empty-value",
																							"use_existing_user": true,
																							"keep_user_at_suite_end": true,
																							"existing_user_password": "non-empty-value",
																							"backend": "non-empty-value",
																							"staticfile_buildpack_name": "non-empty-value",
																							"java_buildpack_name": "non-empty-value",
																							"ruby_buildpack_name": "non-empty-value",
																							"nodejs_buildpack_name": "non-empty-value",
																							"go_buildpack_name": "non-empty-value",
																							"python_buildpack_name": "non-empty-value",
																							"php_buildpack_name": "non-empty-value",
																							"binary_buildpack_name": "non-empty-value",
																							"persistent_app_host": "non-empty-value",
																							"persistent_app_space": "non-empty-value",
																							"persistent_app_org": "non-empty-value",
																							"persistent_app_quota_name": "non-empty-value",
																							"default_timeout": 1,
																							"cf_push_timeout": 1,
																							"long_curl_timeout": 1,
																							"broker_start_timeout": 1
																							}`))

			})
		})

	})

	Describe("writing the integration_confg json file", func() {

		Context("when no env vars are set", func() {

			It("writes the config object to the destination file as json", func() {
				config := configwriter.GenerateConfigFromEnv()
				tempDir := os.TempDir()

				configwriter.WriteConfigToFile(tempDir, config)

				Expect(tempDir + "integration_config.json").To(BeARegularFile())

				file, _ := os.Open(tempDir + "integration_config.json")
				contents, _ := ioutil.ReadFile(file.Name())
				Expect(contents).To(MatchJSON(`{
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

		Context("when env vars are set", func() {
			BeforeEach(func() {
				os.Setenv("CF_API", "cf-api-value")
			})

			AfterEach(func() {
				os.Unsetenv("CF_API")
			})

			It("writes the config object to the destination file as json", func() {
				config := configwriter.GenerateConfigFromEnv()
				tempDir := os.TempDir()

				configwriter.WriteConfigToFile(tempDir, config)

				file, _ := os.Open(tempDir + "integration_config.json")
				contents, _ := ioutil.ReadFile(file.Name())

				Expect(contents).To(MatchJSON(`{
																							"api": "cf-api-value",
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

	})

})
