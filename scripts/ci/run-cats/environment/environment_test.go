package environment_test

import (
	"os"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/environment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	Describe("GetBoolean", func() {
		Context("when the variable is not set", func() {
			It("returns false", func() {
				env := environment.New()
				boolValue, err := env.GetBoolean("MY_ENV_VAR")
				Expect(boolValue).To(BeFalse())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the environment variable is set to the empty string", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns false", func() {
				env := environment.New()
				boolValue, err := env.GetBoolean("MY_ENV_VAR")
				Expect(boolValue).To(BeFalse())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the environment variable is set to the string 'true'", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "true")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns true", func() {
				env := environment.New()
				boolValue, _ := env.GetBoolean("MY_ENV_VAR")
				Expect(boolValue).To(BeTrue())
			})
		})

		Context("when the environment variable is set to a non-boolean", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "not a boolean")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetBoolean("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be a boolean 'true' or 'false'"))
			})
		})

		Context("when the environment variable is set to a non-boolean value that ParseBool would accept", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "T")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetBoolean("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be a boolean 'true' or 'false'"))
			})
		})
	})

	Describe("GetBooleanDefaultToTrue", func() {
		Context("when the variable is not set", func() {
			It("returns true", func() {
				env := environment.New()
				boolValue, err := env.GetBooleanDefaultToTrue("MY_ENV_VAR")
				Expect(boolValue).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the environment variable is set to the empty string", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns true", func() {
				env := environment.New()
				boolValue, err := env.GetBooleanDefaultToTrue("MY_ENV_VAR")
				Expect(boolValue).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the environment variable is set to the string 'true'", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "true")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns true", func() {
				env := environment.New()
				boolValue, _ := env.GetBooleanDefaultToTrue("MY_ENV_VAR")
				Expect(boolValue).To(BeTrue())
			})
		})

		Context("when the environment variable is set to a non-boolean", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "not a boolean")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetBooleanDefaultToTrue("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be a boolean 'true' or 'false'"))
			})
		})

		Context("when the environment variable is set to a non-boolean value that ParseBool would accept", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "T")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetBooleanDefaultToTrue("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be a boolean 'true' or 'false'"))
			})
		})
	})

	Describe("GetInteger", func() {
		Context("when the variable is not set", func() {
			It("returns 0 because it is the default value for integer", func() {
				env := environment.New()
				intValue, err := env.GetInteger("MY_ENV_VAR")
				Expect(err).NotTo(HaveOccurred())
				Expect(intValue).To(Equal(0))
			})
		})

		Context("when the variable is explicitly set to 0", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "0")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetInteger("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be an integer greater than 0"))
			})
		})

		Context("when the variable is not an integer", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "not an integer")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetInteger("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be an integer greater than 0"))
			})
		})
		Context("when the variable is set to a negative integer", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "-1")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetInteger("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be an integer greater than 0"))
			})
		})

		Context("when the variable is set to a strictly positive integer", func() {
			BeforeEach(func() {
				os.Setenv("MY_ENV_VAR", "10")
			})

			AfterEach(func() {
				os.Unsetenv("MY_ENV_VAR")
			})

			It("returns the integer value", func() {
				env := environment.New()
				intValue, err := env.GetInteger("MY_ENV_VAR")
				Expect(intValue).To(Equal(10))
				Expect(err).NotTo(HaveOccurred())
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
			Expect(err.Error()).To(Equal("Invalid environment variable: 'SKIP_SSL_VALIDATION' must be a boolean 'true' or 'false'"))
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
			Expect(err.Error()).To(Equal("Invalid environment variable: 'USE_HTTP' must be a boolean 'true' or 'false'"))
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
			Expect(err.Error()).To(Equal("Invalid environment variable: 'INCLUDE_PRIVILEGED_CONTAINER_SUPPORT' must be a boolean 'true' or 'false'"))
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
			Expect(err.Error()).To(Equal("Invalid environment variable: 'DEFAULT_TIMEOUT_IN_SECONDS' must be an integer greater than 0"))
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
			Expect(err.Error()).To(Equal("Invalid environment variable: 'CF_PUSH_TIMEOUT_IN_SECONDS' must be an integer greater than 0"))
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
			Expect(err.Error()).To(Equal("Invalid environment variable: 'LONG_CURL_TIMEOUT_IN_SECONDS' must be an integer greater than 0"))
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
			Expect(err.Error()).To(Equal("Invalid environment variable: 'BROKER_START_TIMEOUT_IN_SECONDS' must be an integer greater than 0"))
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
				Expect(err.Error()).To(Equal("Invalid environment variable: 'BACKEND' was 'some other backend', but must be 'diego', 'dea', or empty"))
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
})
