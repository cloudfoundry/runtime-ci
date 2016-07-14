package configwriter_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/configwriter"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/environment/fake"
	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configwriter", func() {
	var env *fake.FakeEnvironment
	BeforeEach(func() {
		env = &fake.FakeEnvironment{}
	})

	Context("when env vars are not set", func() {
		It("returns an empty config object", func() {
			configFile, err := configwriter.NewConfigFile("/dir/name", env)

			Expect(err).NotTo(HaveOccurred())
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

		It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to false", func() {
			configFile, err := configwriter.NewConfigFile("", env)
			Expect(err).NotTo(HaveOccurred())
			Expect(configFile).NotTo(BeNil())
			Expect(configFile.Config.UseExistingUser).To(Equal(false))
			Expect(configFile.Config.KeepUserAtSuiteEnd).To(Equal(false))
		})
	})

	Context("When envvars are set", func() {
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
			expectedApi = "api.example.com"
			expectedAdminUser = "admin_user"
			expectedPassword = "admin_password"
			expectedAppsDomain = "apps_domain"
			expectedSkipSslValidation = true
			expectedUseHttp = true
			expectedExistingUser = "existing_user"
			expectedExistingUserPassword = "expected_existing_user_password"
			expectedBackend = "diego"
			expectedPersistedAppHost = "PERSISTENT_APP_HOST"
			expectedPersistedAppSpace = "PERSISTENT_APP_SPACE"
			expectedPersistedAppOrg = "PERSISTENT_APP_ORG"
			expectedPersistedAppQuotaName = "PERSISTENT_APP_QUOTA_NAME"

			expectedDefaultTimeout = 1
			expectedCfPushTimeout = 2
			expectedLongCurlTimeout = 3
			expectedBrokerStartTimeout = 4

			expectedStaticBuildpackName = "STATICFILE_BUILDPACK_NAME"
			expectedJavaBuildpackName = "JAVA_BUILDPACK_NAME"
			expectedRubyBuildpackName = "Ruby_BUILDPACK_NAME"
			expectedNodeJsBuildpackName = "NODEJS_BUILDPACK_NAME"
			expectedGoBuildpackName = "GO_BUILDPACK_NAME"
			expectedPythonBuildpackName = "PYTHON_BUILDPACK_NAME"
			expectedPhpBuildpackName = "PHP_BUILDPACK_NAME"
			expectedBinaryBuildpackName = "BINARY_BUILDPACK_NAME"

			env.GetBooleanReturnsFor("SKIP_SSL_VALIDATION", expectedSkipSslValidation, nil)
			env.GetBooleanReturnsFor("USE_HTTP", expectedUseHttp, nil)

			env.GetStringReturnsFor("CF_API", expectedApi)
			env.GetStringReturnsFor("CF_ADMIN_USER", expectedAdminUser)
			env.GetStringReturnsFor("CF_ADMIN_PASSWORD", expectedPassword)
			env.GetStringReturnsFor("CF_APPS_DOMAIN", expectedAppsDomain)
			env.GetStringReturnsFor("EXISTING_USER", expectedExistingUser)
			env.GetStringReturnsFor("EXISTING_USER_PASSWORD", expectedExistingUserPassword)
			env.GetStringReturnsFor("PERSISTENT_APP_HOST", expectedPersistedAppHost)
			env.GetStringReturnsFor("PERSISTENT_APP_SPACE", expectedPersistedAppSpace)
			env.GetStringReturnsFor("PERSISTENT_APP_ORG", expectedPersistedAppOrg)
			env.GetStringReturnsFor("PERSISTENT_APP_QUOTA_NAME", expectedPersistedAppQuotaName)
			env.GetStringReturnsFor("STATICFILE_BUILDPACK_NAME", expectedStaticBuildpackName)
			env.GetStringReturnsFor("JAVA_BUILDPACK_NAME", expectedJavaBuildpackName)
			env.GetStringReturnsFor("RUBY_BUILDPACK_NAME", expectedRubyBuildpackName)
			env.GetStringReturnsFor("NODEJS_BUILDPACK_NAME", expectedNodeJsBuildpackName)
			env.GetStringReturnsFor("GO_BUILDPACK_NAME", expectedGoBuildpackName)
			env.GetStringReturnsFor("PYTHON_BUILDPACK_NAME", expectedPythonBuildpackName)
			env.GetStringReturnsFor("PHP_BUILDPACK_NAME", expectedPhpBuildpackName)
			env.GetStringReturnsFor("BINARY_BUILDPACK_NAME", expectedBinaryBuildpackName)

			env.GetIntegerReturnsFor("DEFAULT_TIMEOUT_IN_SECONDS", expectedDefaultTimeout, nil)
			env.GetIntegerReturnsFor("CF_PUSH_TIMEOUT_IN_SECONDS", expectedCfPushTimeout, nil)
			env.GetIntegerReturnsFor("LONG_CURL_TIMEOUT_IN_SECONDS", expectedLongCurlTimeout, nil)
			env.GetIntegerReturnsFor("BROKER_START_TIMEOUT_IN_SECONDS", expectedBrokerStartTimeout, nil)

			env.GetBackendReturns(expectedBackend, nil)
		})

		It("Generates a config object with the correct CF env variables set", func() {
			configFile, err := configwriter.NewConfigFile("/some/dir", env)
			Expect(err).NotTo(HaveOccurred())
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

		It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to true", func() {
			configFile, err := configwriter.NewConfigFile("/some/dir", env)
			Expect(err).NotTo(HaveOccurred())
			Expect(configFile.Config.UseExistingUser).To(Equal(true))
			Expect(configFile.Config.KeepUserAtSuiteEnd).To(Equal(true))
		})

		Context("when timeouts are not valid integers", func() {
			DescribeTable("fails fast with the provided error", func(envVarKey string) {
				env.GetIntegerReturnsFor(envVarKey, 0, fmt.Errorf("some error"))
				_, err := configwriter.NewConfigFile("", env)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("some error"))
			},
				Entry("for DEFAULT_TIMEOUT_IN_SECONDS", "DEFAULT_TIMEOUT_IN_SECONDS"),
				Entry("for CF_PUSH_TIMEOUT_IN_SECONDS", "CF_PUSH_TIMEOUT_IN_SECONDS"),
				Entry("for LONG_CURL_TIMEOUT_IN_SECONDS", "LONG_CURL_TIMEOUT_IN_SECONDS"),
				Entry("for BROKER_START_TIMEOUT_IN_SECONDS", "BROKER_START_TIMEOUT_IN_SECONDS"),
			)
		})

		Context("when boolean environment variables are not valid booleans", func() {
			DescribeTable("fails fast with the provided error", func(varName string) {
				env.GetBooleanReturnsFor(varName, false, fmt.Errorf("some boolean error"))
				_, err := configwriter.NewConfigFile("", env)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("some boolean error"))
			},
				Entry("SKIP_SSL_VALIDATION", "SKIP_SSL_VALIDATION"),
				Entry("USE_HTTP", "USE_HTTP"),
			)
		})

		Context("when GetBackend returns an error", func() {
			var expectedError error
			BeforeEach(func() {
				expectedError = fmt.Errorf("some backend error")
				env.GetBackendReturns("", expectedError)
			})

			It("propogates the error", func() {
				_, err := configwriter.NewConfigFile("", env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("some backend error"))
			})
		})
	})

	Describe("marshaling the struct", func() {
		It("does not render optional keys if their values are empty", func() {
			var configJson []byte
			configFile, err := configwriter.NewConfigFile("", env)
			Expect(err).NotTo(HaveOccurred())
			configJson, err = json.Marshal(configFile.Config)
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
				env.GetBooleanReturnsFor("SKIP_SSL_VALIDATION", true, nil)
				env.GetBooleanReturnsFor("USE_HTTP", true, nil)

				env.GetBackendReturns("diego", nil)

				env.GetStringReturnsFor("CF_API", "non-empty-value")
				env.GetStringReturnsFor("CF_ADMIN_USER", "non-empty-value")
				env.GetStringReturnsFor("CF_ADMIN_PASSWORD", "non-empty-value")
				env.GetStringReturnsFor("CF_APPS_DOMAIN", "non-empty-value")
				env.GetStringReturnsFor("EXISTING_USER", "non-empty-value")
				env.GetStringReturnsFor("EXISTING_USER_PASSWORD", "non-empty-value")
				env.GetStringReturnsFor("PERSISTENT_APP_HOST", "non-empty-value")
				env.GetStringReturnsFor("PERSISTENT_APP_SPACE", "non-empty-value")
				env.GetStringReturnsFor("PERSISTENT_APP_ORG", "non-empty-value")
				env.GetStringReturnsFor("PERSISTENT_APP_QUOTA_NAME", "non-empty-value")
				env.GetIntegerReturnsFor("DEFAULT_TIMEOUT_IN_SECONDS", 1, nil)
				env.GetIntegerReturnsFor("CF_PUSH_TIMEOUT_IN_SECONDS", 1, nil)
				env.GetIntegerReturnsFor("LONG_CURL_TIMEOUT_IN_SECONDS", 1, nil)
				env.GetIntegerReturnsFor("BROKER_START_TIMEOUT_IN_SECONDS", 1, nil)
				env.GetStringReturnsFor("STATICFILE_BUILDPACK_NAME", "non-empty-value")
				env.GetStringReturnsFor("JAVA_BUILDPACK_NAME", "non-empty-value")
				env.GetStringReturnsFor("RUBY_BUILDPACK_NAME", "non-empty-value")
				env.GetStringReturnsFor("NODEJS_BUILDPACK_NAME", "non-empty-value")
				env.GetStringReturnsFor("GO_BUILDPACK_NAME", "non-empty-value")
				env.GetStringReturnsFor("PYTHON_BUILDPACK_NAME", "non-empty-value")
				env.GetStringReturnsFor("PHP_BUILDPACK_NAME", "non-empty-value")
				env.GetStringReturnsFor("BINARY_BUILDPACK_NAME", "non-empty-value")
			})

			It("renders the variables in the integration_config", func() {
				var configJson []byte
				configFile, err := configwriter.NewConfigFile("", env)
				Expect(err).NotTo(HaveOccurred())
				configJson, err = json.Marshal(configFile.Config)
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
																							"backend": "diego",
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
		var tempDir string
		var err error

		BeforeEach(func() {
			tempDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			env.GetStringReturnsFor("CF_API", "cf-api-value")
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("writes the marshalled config to the file", func() {
			configFile, err := configwriter.NewConfigFile(tempDir, env)
			Expect(err).NotTo(HaveOccurred())

			var file *os.File
			file, err = configFile.WriteConfigToFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(tempDir + "/integration_config.json").To(BeARegularFile())
			Expect(tempDir + "/integration_config.json").To(Equal(file.Name()))

			contents, err := ioutil.ReadFile(file.Name())
			Expect(err).NotTo(HaveOccurred())
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

		Context("when the destinationDir is invalid", func() {
			It("fails with a nice error", func() {
				configFile, err := configwriter.NewConfigFile("/badpath", env)
				Expect(err).NotTo(HaveOccurred())

				_, err = configFile.WriteConfigToFile()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unable to write integration_config.json to /badpath"))
			})
		})
	})

	Context("when the destinationDir doesn't end with a trailing slash", func() {
		BeforeEach(func() {
			env.GetStringReturnsFor("CF_API", "cf-api-value")
		})

		It("should successfully write integration_config.json", func() {
			configFile, err := configwriter.NewConfigFile("/tmp", env)
			Expect(err).NotTo(HaveOccurred())

			_, err = configFile.WriteConfigToFile()
			Expect(err).NotTo(HaveOccurred())

			var file *os.File
			file, err = os.Open("/tmp/integration_config.json")
			Expect(err).NotTo(HaveOccurred())
			contents, err := ioutil.ReadFile(file.Name())

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
				configFile, err := configwriter.NewConfigFile("/some/path", env)
				Expect(err).NotTo(HaveOccurred())

				configFile.ExportConfigFilePath()

				Expect(os.Getenv("CONFIG")).To(Equal("/some/path/integration_config.json"))
			})
		})
	})
})
