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
			configFile := configwriter.NewConfigFile("/dir/name")
			Expect(configFile).NotTo(BeNil())
			Expect(configFile.Config.Api).To(Equal(""))
			Expect(configFile.Config.AdminUser).To(Equal(""))
			Expect(configFile.Config.AdminPassword).To(Equal(""))
			Expect(configFile.Config.AppsDomain).To(Equal(""))
			Expect(configFile.Config.SkipSslValidation).To(BeFalse())
			Expect(configFile.Config.UseHttp).To(BeFalse())
			Expect(configFile.Config.ExistingUser).To(Equal(""))
			Expect(configFile.Config.ExistingUserPassword).To(Equal(""))
			Expect(configFile.Config.Backend).To(Equal(""))
			Expect(configFile.Config.PersistentAppHost).To(Equal(""))
			Expect(configFile.Config.PersistentAppSpace).To(Equal(""))
			Expect(configFile.Config.PersistentAppOrg).To(Equal(""))
			Expect(configFile.Config.PersistentAppQuotaName).To(Equal(""))
			Expect(configFile.Config.StaticBuildpackName).To(Equal(""))
			Expect(configFile.Config.JavaBuildpackName).To(Equal(""))
			Expect(configFile.Config.RubyBuildpackName).To(Equal(""))
			Expect(configFile.Config.NodeJsBuildpackName).To(Equal(""))
			Expect(configFile.Config.GoBuildpackName).To(Equal(""))
			Expect(configFile.Config.PythonBuildpackName).To(Equal(""))
			Expect(configFile.Config.PhpBuildpackName).To(Equal(""))
			Expect(configFile.Config.BinaryBuildpackName).To(Equal(""))
			Expect(configFile.DestinationDir).To(Equal("/dir/name"))
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
			configFile := configwriter.NewConfigFile("/some/dir")
			Expect(configFile).NotTo(BeNil())
			Expect(configFile.Config.Api).To(Equal(expectedApi))
			Expect(configFile.Config.AdminUser).To(Equal(expectedAdminUser))
			Expect(configFile.Config.AdminPassword).To(Equal(expectedPassword))
			Expect(configFile.Config.AppsDomain).To(Equal(expectedAppsDomain))
			Expect(configFile.Config.SkipSslValidation).To(Equal(expectedSkipSslValidation))
			Expect(configFile.Config.UseHttp).To(Equal(expectedUseHttp))
			Expect(configFile.Config.ExistingUser).To(Equal(expectedExistingUser))
			Expect(configFile.Config.ExistingUserPassword).To(Equal(expectedExistingUserPassword))
			Expect(configFile.Config.Backend).To(Equal(expectedBackend))
			Expect(configFile.Config.PersistentAppHost).To(Equal(expectedPersistedAppHost))
			Expect(configFile.Config.PersistentAppSpace).To(Equal(expectedPersistedAppSpace))
			Expect(configFile.Config.PersistentAppOrg).To(Equal(expectedPersistedAppOrg))
			Expect(configFile.Config.PersistentAppQuotaName).To(Equal(expectedPersistedAppQuotaName))
			Expect(configFile.Config.DefaultTimeout).To(Equal(expectedDefaultTimeout))
			Expect(configFile.Config.CfPushTimeout).To(Equal(expectedCfPushTimeout))
			Expect(configFile.Config.LongCurlTimeout).To(Equal(expectedLongCurlTimeout))
			Expect(configFile.Config.BrokerStartTimeout).To(Equal(expectedBrokerStartTimeout))
			Expect(configFile.Config.StaticBuildpackName).To(Equal(expectedStaticBuildpackName))
			Expect(configFile.Config.JavaBuildpackName).To(Equal(expectedJavaBuildpackName))
			Expect(configFile.Config.RubyBuildpackName).To(Equal(expectedRubyBuildpackName))
			Expect(configFile.Config.NodeJsBuildpackName).To(Equal(expectedNodeJsBuildpackName))
			Expect(configFile.Config.GoBuildpackName).To(Equal(expectedGoBuildpackName))
			Expect(configFile.Config.PythonBuildpackName).To(Equal(expectedPythonBuildpackName))
			Expect(configFile.Config.PhpBuildpackName).To(Equal(expectedPhpBuildpackName))
			Expect(configFile.Config.BinaryBuildpackName).To(Equal(expectedBinaryBuildpackName))
			Expect(configFile.DestinationDir).To(Equal("/some/dir"))
		})

		Context("when 'ExistingUser' is provided", func() {
			BeforeEach(func() {
				os.Setenv("EXISTING_USER", expectedExistingUser)
			})

			It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to true if 'ExistingUser' is provided", func() {
				configFile := configwriter.NewConfigFile("/some/dir")
				Expect(configFile.Config.UseExistingUser).To(Equal(true))
				Expect(configFile.Config.KeepUserAtSuiteEnd).To(Equal(true))
			})
		})

		Context("when 'ExistingUser' is not provided", func() {
			BeforeEach(func() {
				os.Unsetenv("EXISTING_USER")
			})

			It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to false if 'ExistingUser' is not provided", func() {
				configFile := configwriter.NewConfigFile("")
				Expect(configFile).NotTo(BeNil())
				Expect(configFile.Config.UseExistingUser).To(Equal(false))
				Expect(configFile.Config.KeepUserAtSuiteEnd).To(Equal(false))
			})
		})
	})

	Describe("marshaling the struct", func() {
		It("does not render optional keys if their values are empty", func() {
			configJson, err := json.Marshal(configwriter.NewConfigFile("").Config)
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
				configJson, err := json.Marshal(configwriter.NewConfigFile("").Config)
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

	Describe("writing the integration_config.json file", func() {
		Context("when no env vars are set", func() {
			var tempDir string
			var err error

			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err := os.RemoveAll(tempDir)
				Expect(err).NotTo(HaveOccurred())
			})

			It("writes the config object to the destination file as json", func() {
				configFile := configwriter.NewConfigFile(tempDir)

				file, _ := configFile.WriteConfigToFile()

				Expect(tempDir + "integration_config.json").To(BeARegularFile())
				Expect(tempDir + "integration_config.json").To(Equal(file.Name()))

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
			var tempDir string
			var err error

			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				os.Setenv("CF_API", "cf-api-value")
			})

			AfterEach(func() {
				err := os.RemoveAll(tempDir)
				Expect(err).NotTo(HaveOccurred())

				err = os.Unsetenv("CONFIG")
				err = os.Unsetenv("CF_API")
			})

			It("writes the config object to the destination file as json", func() {
				configFile := configwriter.NewConfigFile(tempDir)

				configFile.WriteConfigToFile()

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

		Context("when the destinationDir is invalid", func() {
			It("fails with a nice error", func() {
				configFile := configwriter.NewConfigFile("/badpath")

				_, err := configFile.WriteConfigToFile()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unable to write integration_config.json to /badpath"))
			})
		})
	})

	Describe("exporting the config environment variable", func() {
		Context("when no env vars are set", func() {
			var tempDir string
			var err error

			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err := os.RemoveAll(tempDir)
				Expect(err).NotTo(HaveOccurred())

				os.Unsetenv("CONFIG")
			})

			It("exports the location of the integration_config.json file", func() {
				configFile := configwriter.NewConfigFile("/some/path")

				configFile.ExportConfigFilePath()

				Expect(os.Getenv("CONFIG")).To(Equal("/some/path/integration_config.json"))
			})
		})
	})
})
