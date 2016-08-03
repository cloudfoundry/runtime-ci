package environment_test

import (
	"os"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/environment"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/validationerrors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {

	Describe("Validation", func() {
		var errorMessages string
		var errors validationerrors.Errors
		JustBeforeEach(func() {
			env := environment.New()
			errors = env.Validate()
			errorMessages = errors.Error()
		})

		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
			os.Unsetenv("SKIP_SSL_VALIDATION")
			os.Unsetenv("USE_HTTP")
			os.Unsetenv("BACKEND")
			os.Unsetenv("DEFAULT_TIMEOUT_IN_SECONDS")
			os.Unsetenv("CF_PUSH_TIMEOUT_IN_SECONDS")
			os.Unsetenv("LONG_CURL_TIMEOUT_IN_SECONDS")
			os.Unsetenv("BROKER_START_TIMEOUT_IN_SECONDS")
			os.Unsetenv("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT")
			os.Unsetenv("SKIP_SSO")

			os.Unsetenv("NODES")
			os.Unsetenv("INCLUDE_DIEGO_SSH")
			os.Unsetenv("INCLUDE_V3")
			os.Unsetenv("INCLUDE_DIEGO_DOCKER")
			os.Unsetenv("INCLUDE_BACKEND_COMPATIBILITY")
			os.Unsetenv("INCLUDE_SECURITY_GROUPS")
			os.Unsetenv("INCLUDE_OPERATOR")
			os.Unsetenv("INCLUDE_INTERNET_DEPENDENT")
			os.Unsetenv("INCLUDE_SERVICES")
			os.Unsetenv("INCLUDE_ROUTE_SERVICES")
		})

		Context("when there are errors", func() {
			BeforeEach(func() {
				os.Setenv("CATS_PATH", "fixtures/pass")

				os.Unsetenv("CF_API")
				os.Unsetenv("CF_ADMIN_USER")
				os.Unsetenv("CF_ADMIN_PASSWORD")
				os.Unsetenv("CF_APPS_DOMAIN")
				os.Setenv("SKIP_SSL_VALIDATION", "Righteous")
				os.Setenv("USE_HTTP", "False")
				os.Setenv("BACKEND", "kubernetes")
				os.Setenv("DEFAULT_TIMEOUT_IN_SECONDS", "60s")
				os.Setenv("CF_PUSH_TIMEOUT_IN_SECONDS", "120 years")
				os.Setenv("LONG_CURL_TIMEOUT_IN_SECONDS", "180 days")
				os.Setenv("BROKER_START_TIMEOUT_IN_SECONDS", "240 mins")
				os.Setenv("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT", "true\n")
				os.Setenv("SKIP_SSO", "falsey")

				os.Setenv("NODES", "five")
				os.Setenv("INCLUDE_DIEGO_SSH", "1")
				os.Setenv("INCLUDE_V3", "0")
				os.Setenv("INCLUDE_DIEGO_DOCKER", "diego")
				os.Setenv("INCLUDE_BACKEND_COMPATIBILITY", "no")
				os.Setenv("INCLUDE_SECURITY_GROUPS", "yes")
				os.Setenv("INCLUDE_OPERATOR", "F")
				os.Setenv("INCLUDE_INTERNET_DEPENDENT", "Falz")
				os.Setenv("INCLUDE_SERVICES", "troo")
				os.Setenv("INCLUDE_ROUTE_SERVICES", "truce")
			})

			It("contains all the error messages", func() {
				Expect(errorMessages).To(And(
					ContainSubstring(`* Missing required environment variables:
    CF_API
    CF_ADMIN_USER
    CF_ADMIN_PASSWORD
    CF_APPS_DOMAIN`),
					ContainSubstring("* Invalid environment variable: 'SKIP_SSL_VALIDATION' must be a boolean 'true' or 'false' but was set to 'Righteous'"),
					ContainSubstring("* Invalid environment variable: 'USE_HTTP' must be a boolean 'true' or 'false' but was set to 'False'"),
					ContainSubstring("* Invalid environment variable: 'BACKEND' must be 'diego', 'dea', or empty but was set to 'kubernetes'"),
					ContainSubstring("* Invalid environment variable: 'DEFAULT_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to '60s'"),
					ContainSubstring("* Invalid environment variable: 'CF_PUSH_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to '120 years'"),
					ContainSubstring("* Invalid environment variable: 'LONG_CURL_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to '180 days'"),
					ContainSubstring("* Invalid environment variable: 'BROKER_START_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to '240 mins'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_PRIVILEGED_CONTAINER_SUPPORT' must be a boolean 'true' or 'false' but was set to 'true\n'"),
					ContainSubstring("* Invalid environment variable: 'SKIP_SSO' must be a boolean 'true' or 'false' but was set to 'falsey'"),
					ContainSubstring("* Invalid environment variable: 'NODES' must be an integer greater than 0 but was set to 'five'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_DIEGO_SSH' must be a boolean 'true' or 'false' but was set to '1'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_V3' must be a boolean 'true' or 'false' but was set to '0'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_DIEGO_DOCKER' must be a boolean 'true' or 'false' but was set to 'diego'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_BACKEND_COMPATIBILITY' must be a boolean 'true' or 'false' but was set to 'no'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_SECURITY_GROUPS' must be a boolean 'true' or 'false' but was set to 'yes'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_OPERATOR' must be a boolean 'true' or 'false' but was set to 'F'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_INTERNET_DEPENDENT' must be a boolean 'true' or 'false' but was set to 'Falz'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_SERVICES' must be a boolean 'true' or 'false' but was set to 'troo'"),
					ContainSubstring("* Invalid environment variable: 'INCLUDE_ROUTE_SERVICES' must be a boolean 'true' or 'false' but was set to 'truce'"),
				))
			})
		})

		Context("when there are no errors", func() {
			BeforeEach(func() {
				os.Setenv("CF_API", "blah")
				os.Setenv("CF_ADMIN_USER", "blah2")
				os.Setenv("CF_ADMIN_PASSWORD", "oogieboogie")
				os.Setenv("CF_APPS_DOMAIN", "example.com")
			})

			It("should be empty", func() {
				Expect(errors.Empty()).To(BeTrue())
			})
		})
	})
	Describe("GetSkipSSLValidation", func() {
		AfterEach(func() {
			os.Unsetenv("SKIP_SSL_VALIDATION")
		})

		It("returns a boolean when set properly", func() {
			env := environment.New()
			os.Setenv("SKIP_SSL_VALIDATION", "true")
			result, _ := env.GetSkipSSLValidation()
			Expect(result).To(BeTrue())
		})

		It("returns a default of false", func() {
			env := environment.New()
			result, _ := env.GetSkipSSLValidation()
			Expect(result).To(BeFalse())
		})

		It("returns an error when it is set wrong", func() {
			env := environment.New()
			os.Setenv("SKIP_SSL_VALIDATION", "blah blah blah")
			_, err := env.GetSkipSSLValidation()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'SKIP_SSL_VALIDATION' must be a boolean 'true' or 'false' but was set to 'blah blah blah'"))
		})
	})

	Describe("GetUseHTTP", func() {
		AfterEach(func() {
			os.Unsetenv("USE_HTTP")
		})

		It("returns a boolean when set properly", func() {
			env := environment.New()
			os.Setenv("USE_HTTP", "true")
			result, _ := env.GetUseHTTP()
			Expect(result).To(BeTrue())
		})

		It("returns a default of false", func() {
			env := environment.New()
			result, _ := env.GetUseHTTP()
			Expect(result).To(BeFalse())
		})

		It("returns an error when it is set wrong", func() {
			env := environment.New()
			os.Setenv("USE_HTTP", "blah blah blah")
			_, err := env.GetUseHTTP()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'USE_HTTP' must be a boolean 'true' or 'false' but was set to 'blah blah blah'"))
		})
	})

	Describe("GetIncludePrivilegedContainerSupport", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT")
		})

		It("returns a boolean when set properly", func() {
			env := environment.New()
			os.Setenv("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT", "true")
			result, _ := env.GetIncludePrivilegedContainerSupport()
			Expect(result).To(BeTrue())
		})

		It("returns a default of false", func() {
			env := environment.New()
			result, _ := env.GetIncludePrivilegedContainerSupport()
			Expect(result).To(BeFalse())
		})

		It("returns an error when it is set wrong", func() {
			env := environment.New()
			os.Setenv("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT", "blah blah blah")
			_, err := env.GetIncludePrivilegedContainerSupport()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_PRIVILEGED_CONTAINER_SUPPORT' must be a boolean 'true' or 'false' but was set to 'blah blah blah'"))
		})
	})

	Describe("GetDefaultTimeoutInSeconds", func() {
		AfterEach(func() {
			os.Unsetenv("DEFAULT_TIMEOUT_IN_SECONDS")
		})

		It("returns an integer when set properly", func() {
			env := environment.New()
			os.Setenv("DEFAULT_TIMEOUT_IN_SECONDS", "2")
			result, _ := env.GetDefaultTimeoutInSeconds()
			Expect(result).To(Equal(2))
		})

		It("returns a default of 0", func() {
			env := environment.New()
			result, _ := env.GetDefaultTimeoutInSeconds()
			Expect(result).To(Equal(0))
		})

		It("returns an error when it is set wrong", func() {
			env := environment.New()
			os.Setenv("DEFAULT_TIMEOUT_IN_SECONDS", "blah blah blah")
			_, err := env.GetDefaultTimeoutInSeconds()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'DEFAULT_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to 'blah blah blah'"))
		})
	})

	Describe("GetCFPushTimeoutInSeconds", func() {
		AfterEach(func() {
			os.Unsetenv("CF_PUSH_TIMEOUT_IN_SECONDS")
		})

		It("returns an integer when set properly", func() {
			env := environment.New()
			os.Setenv("CF_PUSH_TIMEOUT_IN_SECONDS", "2")
			result, _ := env.GetCFPushTimeoutInSeconds()
			Expect(result).To(Equal(2))
		})

		It("returns a default of 0", func() {
			env := environment.New()
			result, _ := env.GetCFPushTimeoutInSeconds()
			Expect(result).To(Equal(0))
		})

		It("returns an error when it is set wrong", func() {
			env := environment.New()
			os.Setenv("CF_PUSH_TIMEOUT_IN_SECONDS", "blah blah blah")
			_, err := env.GetCFPushTimeoutInSeconds()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'CF_PUSH_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to 'blah blah blah'"))
		})
	})

	Describe("GetLongCurlTimeoutInSeconds", func() {
		AfterEach(func() {
			os.Unsetenv("LONG_CURL_TIMEOUT_IN_SECONDS")
		})

		It("returns an integer when set properly", func() {
			env := environment.New()
			os.Setenv("LONG_CURL_TIMEOUT_IN_SECONDS", "2")
			result, _ := env.GetLongCurlTimeoutInSeconds()
			Expect(result).To(Equal(2))
		})

		It("returns a default of 0", func() {
			env := environment.New()
			result, _ := env.GetLongCurlTimeoutInSeconds()
			Expect(result).To(Equal(0))
		})

		It("returns an error when it is set wrong", func() {
			env := environment.New()
			os.Setenv("LONG_CURL_TIMEOUT_IN_SECONDS", "blah blah blah")
			_, err := env.GetLongCurlTimeoutInSeconds()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'LONG_CURL_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to 'blah blah blah'"))
		})
	})

	Describe("GetBrokerStartTimeoutInSeconds", func() {
		AfterEach(func() {
			os.Unsetenv("BROKER_START_TIMEOUT_IN_SECONDS")
		})

		It("returns an integer when set properly", func() {
			env := environment.New()
			os.Setenv("BROKER_START_TIMEOUT_IN_SECONDS", "2")
			result, _ := env.GetBrokerStartTimeoutInSeconds()
			Expect(result).To(Equal(2))
		})

		It("returns a default of 0", func() {
			env := environment.New()
			result, _ := env.GetBrokerStartTimeoutInSeconds()
			Expect(result).To(Equal(0))
		})

		It("returns an error when it is set wrong", func() {
			env := environment.New()
			os.Setenv("BROKER_START_TIMEOUT_IN_SECONDS", "blah blah blah")
			_, err := env.GetBrokerStartTimeoutInSeconds()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'BROKER_START_TIMEOUT_IN_SECONDS' must be an integer greater than 0 but was set to 'blah blah blah'"))
		})
	})

	Describe("GetBackend", func() {
		Context("when the value is empty", func() {
			It("returns the empty string", func() {
				env := environment.New()
				backend, err := env.GetBackend()
				Expect(err).NotTo(HaveOccurred())
				Expect(backend).To(Equal(""))
			})
		})

		Context("when the value is 'diego'", func() {
			BeforeEach(func() {
				os.Setenv("BACKEND", "diego")
			})

			It("returns 'diego'", func() {
				env := environment.New()
				backend, err := env.GetBackend()
				Expect(err).NotTo(HaveOccurred())
				Expect(backend).To(Equal("diego"))
			})
		})

		Context("when the value is 'dea'", func() {
			BeforeEach(func() {
				os.Setenv("BACKEND", "dea")
			})

			It("returns 'dea'", func() {
				env := environment.New()
				backend, err := env.GetBackend()
				Expect(err).NotTo(HaveOccurred())
				Expect(backend).To(Equal("dea"))
			})
		})

		Context("when the value is anything else", func() {
			BeforeEach(func() {
				os.Setenv("BACKEND", "some other backend")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetBackend()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("* Invalid environment variable: 'BACKEND' must be 'diego', 'dea', or empty but was set to 'some other backend'"))
			})
		})
	})

	Describe("GetCFAPI", func() {
		AfterEach(func() {
			os.Unsetenv("CF_API")
		})

		It("Returns the value set in the CF_API variable", func() {
			expectedResult := "boggles"
			os.Setenv("CF_API", expectedResult)
			env := environment.New()
			result := env.GetCFAPI()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetCFAdminUser", func() {
		AfterEach(func() {
			os.Unsetenv("CF_ADMIN_USER")
		})

		It("Returns the value set in the CF_ADMIN_USER variable", func() {
			expectedResult := "boggles"
			os.Setenv("CF_ADMIN_USER", expectedResult)
			env := environment.New()
			result := env.GetCFAdminUser()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetCFAdminPassword", func() {
		AfterEach(func() {
			os.Unsetenv("CF_ADMIN_PASSWORD")
		})

		It("Returns the value set in the CF_ADMIN_PASSWORD variable", func() {
			expectedResult := "boggles"
			os.Setenv("CF_ADMIN_PASSWORD", expectedResult)
			env := environment.New()
			result := env.GetCFAdminPassword()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetCFAppsDomain", func() {
		AfterEach(func() {
			os.Unsetenv("CF_APPS_DOMAIN")
		})

		It("Returns the value set in the CF_APPS_DOMAIN variable", func() {
			expectedResult := "boggles"
			os.Setenv("CF_APPS_DOMAIN", expectedResult)
			env := environment.New()
			result := env.GetCFAppsDomain()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetExistingUser", func() {
		AfterEach(func() {
			os.Unsetenv("EXISTING_USER")
		})

		It("Returns the value set in the EXISTING_USER variable", func() {
			expectedResult := "boggles"
			os.Setenv("EXISTING_USER", expectedResult)
			env := environment.New()
			result := env.GetExistingUser()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("UseExistingUser", func() {
		AfterEach(func() {
			os.Unsetenv("EXISTING_USER")
		})

		It("Returns true if EXISTING_USER is set", func() {
			os.Setenv("EXISTING_USER", "anything at all")
			env := environment.New()
			result := env.UseExistingUser()
			Expect(result).To(BeTrue())
		})

		It("Returns false if EXISTING_USER is not set or empty", func() {
			env := environment.New()
			result := env.UseExistingUser()
			Expect(result).To(BeFalse())
		})
	})

	Describe("KeepUserAtSuiteEnd", func() {
		AfterEach(func() {
			os.Unsetenv("EXISTING_USER")
		})

		It("Returns true if EXISTING_USER is set", func() {
			os.Setenv("EXISTING_USER", "anything at all")
			env := environment.New()
			result := env.KeepUserAtSuiteEnd()
			Expect(result).To(BeTrue())
		})

		It("Returns false if EXISTING_USER is not set or empty", func() {
			env := environment.New()
			result := env.UseExistingUser()
			Expect(result).To(BeFalse())
		})
	})

	Describe("GetExistingUserPassword", func() {
		AfterEach(func() {
			os.Unsetenv("EXISTING_USER_PASSWORD")
		})

		It("Returns the value set in the EXISTING_USER_PASSWORD variable", func() {
			expectedResult := "boggles"
			os.Setenv("EXISTING_USER_PASSWORD", expectedResult)
			env := environment.New()
			result := env.GetExistingUserPassword()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetStaticBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("JAVA_BUILDPACK_NAME")
		})

		It("Returns the value set in the STATICFILE_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("STATICFILE_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetStaticBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetJavaBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("JAVA_BUILDPACK_NAME")
		})

		It("Returns the value set in the JAVA_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("JAVA_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetJavaBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetRubyBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("RUBY_BUILDPACK_NAME")
		})

		It("Returns the value set in the RUBY_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("RUBY_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetRubyBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetNodeJSBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("NODEJS_BUILDPACK_NAME")
		})

		It("Returns the value set in the NODEJS_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("NODEJS_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetNodeJSBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetGoBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("GO_BUILDPACK_NAME")
		})

		It("Returns the value set in the GO_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("GO_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetGoBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetPythonBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("PYTHON_BUILDPACK_NAME")
		})

		It("Returns the value set in the PYTHON_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("PYTHON_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetPythonBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetPHPBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("PHP_BUILDPACK_NAME")
		})

		It("Returns the value set in the PHP_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("PHP_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetPHPBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetBinaryBuildpackName", func() {
		AfterEach(func() {
			os.Unsetenv("BINARY_BUILDPACK_NAME")
		})

		It("Returns the value set in the BINARY_BUILDPACK_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("BINARY_BUILDPACK_NAME", expectedResult)
			env := environment.New()
			result := env.GetBinaryBuildpackName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetPersistentAppHost", func() {
		AfterEach(func() {
			os.Unsetenv("PERSISTENT_APP_HOST")
		})

		It("Returns the value set in the PERSISTENT_APP_HOST variable", func() {
			expectedResult := "boggles"
			os.Setenv("PERSISTENT_APP_HOST", expectedResult)
			env := environment.New()
			result := env.GetPersistentAppHost()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetPersistentAppSpace", func() {
		AfterEach(func() {
			os.Unsetenv("PERSISTENT_APP_SPACE")
		})

		It("Returns the value set in the PERSISTENT_APP_SPACE variable", func() {
			expectedResult := "boggles"
			os.Setenv("PERSISTENT_APP_SPACE", expectedResult)
			env := environment.New()
			result := env.GetPersistentAppSpace()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetPersistentAppOrg", func() {
		AfterEach(func() {
			os.Unsetenv("PERSISTENT_APP_ORG")
		})

		It("Returns the value set in the PERSISTENT_APP_ORG variable", func() {
			expectedResult := "boggles"
			os.Setenv("PERSISTENT_APP_ORG", expectedResult)
			env := environment.New()
			result := env.GetPersistentAppOrg()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetPersistentAppQuotaName", func() {
		AfterEach(func() {
			os.Unsetenv("PERSISTENT_APP_QUOTA_NAME")
		})

		It("Returns the value set in the PERSISTENT_APP_QUOTA_NAME variable", func() {
			expectedResult := "boggles"
			os.Setenv("PERSISTENT_APP_QUOTA_NAME", expectedResult)
			env := environment.New()
			result := env.GetPersistentAppQuotaName()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("GetSkipDiegoSSH", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_DIEGO_SSH")
		})

		It("Returns the string 'ssh' if it should be skipped", func() {
			os.Setenv("INCLUDE_DIEGO_SSH", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipDiegoSSH()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("ssh"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_DIEGO_SSH", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipDiegoSSH()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_DIEGO_SSH")
			env := environment.New()
			skipFlag, err := env.GetSkipDiegoSSH()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("ssh"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_DIEGO_SSH", "falsey")
			env := environment.New()
			_, err := env.GetSkipDiegoSSH()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_DIEGO_SSH' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipV3", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_V3")
		})

		It("Returns the string 'v3' if it should be skipped", func() {
			os.Setenv("INCLUDE_V3", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipV3()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("v3"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_V3", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipV3()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_V3")
			env := environment.New()
			skipFlag, err := env.GetSkipV3()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("v3"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_V3", "falsey")
			env := environment.New()
			_, err := env.GetSkipV3()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_V3' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipSSO", func() {
		AfterEach(func() {
			os.Unsetenv("SKIP_SSO")
		})

		It("Returns the string 'SSO' if it should be skipped", func() {
			os.Setenv("SKIP_SSO", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipSSO()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("SSO"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("SKIP_SSO", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipSSO()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("SKIP_SSO")
			env := environment.New()
			skipFlag, err := env.GetSkipSSO()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("SSO"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("SKIP_SSO", "falsey")
			env := environment.New()
			_, err := env.GetSkipSSO()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'SKIP_SSO' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipDiegoDocker", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_DIEGO_DOCKER")
		})

		It("Returns the string 'docker' if it should be skipped", func() {
			os.Setenv("INCLUDE_DIEGO_DOCKER", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipDiegoDocker()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("docker"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_DIEGO_DOCKER", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipDiegoDocker()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_DIEGO_DOCKER")
			env := environment.New()
			skipFlag, err := env.GetSkipDiegoDocker()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("docker"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_DIEGO_DOCKER", "falsey")
			env := environment.New()
			_, err := env.GetSkipDiegoDocker()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_DIEGO_DOCKER' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipBackendCompatibility", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_BACKEND_COMPATIBILITY")
		})

		It("Returns the string 'backend_compatibility' if it should be skipped", func() {
			os.Setenv("INCLUDE_BACKEND_COMPATIBILITY", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipBackendCompatibility()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("backend_compatibility"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_BACKEND_COMPATIBILITY", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipBackendCompatibility()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_BACKEND_COMPATIBILITY")
			env := environment.New()
			skipFlag, err := env.GetSkipBackendCompatibility()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("backend_compatibility"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_BACKEND_COMPATIBILITY", "falsey")
			env := environment.New()
			_, err := env.GetSkipBackendCompatibility()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_BACKEND_COMPATIBILITY' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipSecurityGroups", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_SECURITY_GROUPS")
		})

		It("Returns the string 'security_groups' if it should be skipped", func() {
			os.Setenv("INCLUDE_SECURITY_GROUPS", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipSecurityGroups()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("security_groups"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_SECURITY_GROUPS", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipSecurityGroups()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_SECURITY_GROUPS")
			env := environment.New()
			skipFlag, err := env.GetSkipSecurityGroups()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("security_groups"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_SECURITY_GROUPS", "falsey")
			env := environment.New()
			_, err := env.GetSkipSecurityGroups()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_SECURITY_GROUPS' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipOperator", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_OPERATOR")
		})

		It("Returns the string 'operator' if it should be skipped", func() {
			os.Setenv("INCLUDE_OPERATOR", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipOperator()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("operator"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_OPERATOR", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipOperator()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_OPERATOR")
			env := environment.New()
			skipFlag, err := env.GetSkipOperator()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("operator"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_OPERATOR", "falsey")
			env := environment.New()
			_, err := env.GetSkipOperator()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_OPERATOR' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipInternetDependent", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_INTERNET_DEPENDENT")
		})

		It("Returns the string 'internet_dependent' if it should be skipped", func() {
			os.Setenv("INCLUDE_INTERNET_DEPENDENT", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipInternetDependent()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("internet_dependent"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_INTERNET_DEPENDENT", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipInternetDependent()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_INTERNET_DEPENDENT")
			env := environment.New()
			skipFlag, err := env.GetSkipInternetDependent()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("internet_dependent"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_INTERNET_DEPENDENT", "falsey")
			env := environment.New()
			_, err := env.GetSkipInternetDependent()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_INTERNET_DEPENDENT' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipservices", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_SERVICES")
		})

		It("Returns the string 'services' if it should be skipped", func() {
			os.Setenv("INCLUDE_SERVICES", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("services"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_SERVICES", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_SERVICES")
			env := environment.New()
			skipFlag, err := env.GetSkipServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("services"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_SERVICES", "falsey")
			env := environment.New()
			_, err := env.GetSkipServices()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_SERVICES' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetSkipRouteServices", func() {
		AfterEach(func() {
			os.Unsetenv("INCLUDE_ROUTE_SERVICES")
		})

		It("Returns the string 'route_services' if it should be skipped", func() {
			os.Setenv("INCLUDE_ROUTE_SERVICES", "false")
			env := environment.New()
			skipFlag, err := env.GetSkipRouteServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("route_services"))
		})

		It("Returns an empty string if it should not be skipped", func() {
			os.Setenv("INCLUDE_ROUTE_SERVICES", "true")
			env := environment.New()
			skipFlag, err := env.GetSkipRouteServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal(""))
		})

		It("Defaults to skip", func() {
			os.Unsetenv("INCLUDE_ROUTE_SERVICES")
			env := environment.New()
			skipFlag, err := env.GetSkipRouteServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(skipFlag).To(Equal("route_services"))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("INCLUDE_ROUTE_SERVICES", "falsey")
			env := environment.New()
			_, err := env.GetSkipRouteServices()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'INCLUDE_ROUTE_SERVICES' must be a boolean 'true' or 'false' but was set to 'falsey'"))
		})
	})

	Describe("GetBackend", func() {
		AfterEach(func() {
			os.Unsetenv("BACKEND")
		})

		It("Returns the string 'diego' if it set to diego", func() {
			os.Setenv("BACKEND", "diego")
			env := environment.New()
			backend, err := env.GetBackend()
			Expect(err).NotTo(HaveOccurred())
			Expect(backend).To(Equal("diego"))
		})

		It("Returns the string 'dea' if it set to dea", func() {
			os.Setenv("BACKEND", "dea")
			env := environment.New()
			backend, err := env.GetBackend()
			Expect(err).NotTo(HaveOccurred())
			Expect(backend).To(Equal("dea"))
		})

		It("Returns an empty string if it is unset", func() {
			env := environment.New()
			backend, err := env.GetBackend()
			Expect(err).NotTo(HaveOccurred())
			Expect(backend).To(Equal(""))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("BACKEND", "what-is-backend?")
			env := environment.New()
			_, err := env.GetBackend()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'BACKEND' must be 'diego', 'dea', or empty but was set to 'what-is-backend?'"))
		})
	})

	Describe("GetCatsPath", func() {
		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
		})

		It("Returns the value of CATS_PATH ", func() {
			expectedCatsPath := "abcd"
			os.Setenv("CATS_PATH", expectedCatsPath)
			env := environment.New()
			Expect(env.GetCatsPath()).To(Equal(expectedCatsPath))
		})

		It("Has a reasonable default CATS_PATH", func() {
			env := environment.New()
			Expect(env.GetCatsPath()).To(Equal(env.GetGoPath() + "/src/github.com/cloudfoundry/cf-acceptance-tests"))
		})
	})

	Describe("GetNodes", func() {
		AfterEach(func() {
			os.Unsetenv("NODES")
		})

		It("Returns the integer node count", func() {
			os.Setenv("NODES", "3")
			env := environment.New()
			nodes, err := env.GetNodes()
			Expect(err).NotTo(HaveOccurred())
			Expect(nodes).To(Equal(3))
		})

		It("Returns an error if there was an invalid value", func() {
			os.Setenv("NODES", "bogus")
			env := environment.New()
			_, err := env.GetNodes()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("* Invalid environment variable: 'NODES' must be an integer greater than 0 but was set to 'bogus'"))
		})
	})

	Describe("GetGoPath", func() {
		It("Returns the GOPATH", func() {
			env := environment.New()
			Expect(env.GetGoPath()).To(Equal(os.Getenv("GOPATH")))
		})
	})
})
