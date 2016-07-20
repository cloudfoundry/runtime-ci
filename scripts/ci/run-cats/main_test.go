package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	configJsonPath := os.Getenv("PWD") + "/integration_config.json"

	AfterEach(func() {
		os.Remove(configJsonPath)
	})

	Context("when only some required envvars are set", func() {
		BeforeEach(func() {
			os.Setenv("CF_API", "cf_api_value")
		})
		AfterEach(func() {
			os.Unsetenv("CF_API")
		})

		It("Exits 1 and prints only the missing 'required' env vars", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(1))
			Eventually(session.Err, 30).Should(gbytes.Say(`Missing required environment variables:
CF_ADMIN_USER
CF_ADMIN_PASSWORD
CF_APPS_DOMAIN`,
			))
			Expect(configJsonPath).NotTo(BeARegularFile())

		})
	})

	Context("When all required env vars are set", func() {
		BeforeEach(func() {
			//CATS_PATH isn't required, but we need it for test setup
			os.Setenv("CATS_PATH", "fixtures/pass")

			os.Setenv("CF_API", "non-empty-value")
			os.Setenv("CF_ADMIN_USER", "non-empty-value")
			os.Setenv("CF_ADMIN_PASSWORD", "non-empty-value")
			os.Setenv("CF_APPS_DOMAIN", "non-empty-value")
		})
		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
		})

		It("Sets CONFIG envvar to the path for the generated config file", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))

			Expect(string(session.Out.Contents())).To(ContainSubstring("CONFIG=" + configJsonPath))
		})

		It("Writes a config file for CATs to use", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Expect(configJsonPath).To(BeARegularFile())

			configBytes, err := ioutil.ReadFile(configJsonPath)
			Expect(err).NotTo(HaveOccurred())

			var config struct {
				Api               string `json:"api"`
				AdminUser         string `json:"admin_user"`
				AdminPassword     string `json:"admin_password"`
				AppsDomain        string `json:"apps_domain"`
				SkipSSLValidation bool   `json:"skip_ssl_validation"`
				UseHTTP           bool   `json:"use_http"`
			}

			err = json.Unmarshal(configBytes, &config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Api).To(Equal("non-empty-value"))
			Expect(config.AdminUser).To(Equal("non-empty-value"))
			Expect(config.AdminPassword).To(Equal("non-empty-value"))
			Expect(config.AppsDomain).To(Equal("non-empty-value"))
			Expect(config.SkipSSLValidation).To(BeFalse())
			Expect(config.UseHTTP).To(BeFalse())
		})

		It("Executes the command to run CATs, excluding configarable suites and SSO", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Eventually(session.Out, 30).Should(gbytes.Say(
				`bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes=2 -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 -skip=SSO\|NO_DEA_SUPPORT\|NO_DIEGO_SUPPORT -keepGoing`,
			))
		})

		Context("When invalid TIMEOUT env vars are set", func() {
			BeforeEach(func() {
				os.Setenv("DEFAULT_TIMEOUT_IN_SECONDS", "not a timeout")
			})
			AfterEach(func() {
				os.Unsetenv("DEFAULT_TIMEOUT_IN_SECONDS")
			})

			It("Should exit with an appropriate error message and exit code", func() {
				command := exec.Command(binPath)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 30).Should(gexec.Exit(1))
				Eventually(session.Err, 30).Should(gbytes.Say(`Invalid environment variable: 'DEFAULT_TIMEOUT_IN_SECONDS' must be an integer greater than 0`))
				Expect(configJsonPath).NotTo(BeARegularFile())
			})
		})

		Context("When invalid boolean env vars are set", func() {
			BeforeEach(func() {
				os.Setenv("USE_HTTP", "this is not a boolean this is only a tribute")
			})

			AfterEach(func() {
				os.Unsetenv("USE_HTTP")
			})

			It("Should exit with an appropriate error message and exit code", func() {
				command := exec.Command(binPath)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 30).Should(gexec.Exit(1))
				Eventually(session.Err, 30).Should(gbytes.Say(`Invalid environment variable: 'USE_HTTP' must be a boolean 'true' or 'false'`))
				Expect(configJsonPath).NotTo(BeARegularFile())
			})
		})

		Describe("when INCLUDE_* env vars are not proper booleans", func() {
			BeforeEach(func() {
				os.Setenv("INCLUDE_V3", "this is not a boolean this is only a tribute")
			})

			AfterEach(func() {
				os.Unsetenv("INCLUDE_V3")
			})

			It("Should exit with an appropriate error message and exit code", func() {
				command := exec.Command(binPath)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 30).Should(gexec.Exit(1))
				Eventually(session.Err, 30).Should(gbytes.Say(`Invalid environment variable: 'INCLUDE_V3' must be a boolean 'true' or 'false'`))
			})

		})

		Context("When bin/test fails", func() {
			BeforeEach(func() {
				os.Setenv("CATS_PATH", "fixtures/fail")
			})

			AfterEach(func() {
				os.Unsetenv("CATS_PATH")
			})

			It("Returns the same error code as the failing bin/test", func() {
				command := exec.Command(binPath)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 30).Should(gexec.Exit(1))
			})
		})
	})

	Context("When no required env vars are set", func() {
		BeforeEach(func() {
			//CATS_PATH isn't required, but we need it for test setup
			os.Setenv("CATS_PATH", "fixtures/pass")

			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
		})
		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
		})

		It("Exits 1 and prints a list of all 'required' env vars", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(1))
			Eventually(session.Err, 30).Should(gbytes.Say(`Missing required environment variables:
CF_API
CF_ADMIN_USER
CF_ADMIN_PASSWORD
CF_APPS_DOMAIN`,
			))

			Expect(configJsonPath).NotTo(BeARegularFile())
		})

		It("Doesn't write a config file for CATs to use", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 3).Should(gexec.Exit(1))
			Expect(configJsonPath).NotTo(BeARegularFile())
		})

		It("Doesn't execute the command to run CATs, excluding configurable suites and SSO", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 3).Should(gexec.Exit(1))
			Eventually(session.Out, 3).ShouldNot(gbytes.Say(
				`bin/test`,
			))
		})
		Context("and an optional env var is set to an invalid value", func() {
			BeforeEach(func() {
				os.Setenv("SKIP_SSL_VALIDATION", "yes, I would like that please.")
			})
			AfterEach(func() {
				os.Unsetenv("SKIP_SSL_VALIDATION")
			})

			It("displays both the list of missing required vars and the error for the invalid var", func() {
				command := exec.Command(binPath)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 30).Should(gexec.Exit(1))
				Eventually(session.Err, 30).Should(gbytes.Say(`Missing required environment variables:
CF_API
CF_ADMIN_USER
CF_ADMIN_PASSWORD
CF_APPS_DOMAIN`,
				))
				Eventually(session.Err, 3).Should(gbytes.Say(`Invalid environment variable: 'SKIP_SSL_VALIDATION' must be a boolean 'true' or 'false'`))

				Expect(configJsonPath).NotTo(BeARegularFile())
			})
		})
	})

	Context("When all supported env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("CATS_PATH", "fixtures/pass")

			os.Setenv("CF_API", "api.example.com")
			os.Setenv("CF_ADMIN_USER", "admin-username")
			os.Setenv("CF_ADMIN_PASSWORD", "admin-password")
			os.Setenv("CF_APPS_DOMAIN", "apps.example.com")
			os.Setenv("SKIP_SSL_VALIDATION", "true")
			os.Setenv("USE_HTTP", "false")
			os.Setenv("EXISTING_USER", "existing-cats-user")
			os.Setenv("EXISTING_USER_PASSWORD", "existing-cats-user-password")
			os.Setenv("BACKEND", "diego")
			os.Setenv("PERSISTENT_APP_HOST", "cats-app-host")
			os.Setenv("PERSISTENT_APP_SPACE", "cats-app-space")
			os.Setenv("PERSISTENT_APP_ORG", "cats-app-org")
			os.Setenv("PERSISTENT_APP_QUOTA_NAME", "cats-app-quota")
			os.Setenv("DEFAULT_TIMEOUT_IN_SECONDS", "60")
			os.Setenv("CF_PUSH_TIMEOUT_IN_SECONDS", "120")
			os.Setenv("LONG_CURL_TIMEOUT_IN_SECONDS", "180")
			os.Setenv("BROKER_START_TIMEOUT_IN_SECONDS", "240")
			os.Setenv("STATICFILE_BUILDPACK_NAME", "static-buildpack")
			os.Setenv("JAVA_BUILDPACK_NAME", "java-buildpack")
			os.Setenv("RUBY_BUILDPACK_NAME", "ruby-buildpack")
			os.Setenv("NODEJS_BUILDPACK_NAME", "node-buildpack")
			os.Setenv("GO_BUILDPACK_NAME", "go-buildpack")
			os.Setenv("PYTHON_BUILDPACK_NAME", "python-buildpack")
			os.Setenv("PHP_BUILDPACK_NAME", "php-buildpack")
			os.Setenv("BINARY_BUILDPACK_NAME", "binary-buildpack")
			os.Setenv("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT", "true")
			os.Setenv("SKIP_SSO", "false")

			os.Setenv("NODES", "5")
			os.Setenv("INCLUDE_DIEGO_SSH", "true")
			os.Setenv("INCLUDE_V3", "true")
			os.Setenv("INCLUDE_DIEGO_DOCKER", "true")
			os.Setenv("INCLUDE_BACKEND_COMPATIBILITY", "true")
			os.Setenv("INCLUDE_SECURITY_GROUPS", "true")
			os.Setenv("INCLUDE_LOGGING", "true")
			os.Setenv("INCLUDE_OPERATOR", "true")
			os.Setenv("INCLUDE_INTERNET_DEPENDENT", "true")
			os.Setenv("INCLUDE_SERVICES", "true")
			os.Setenv("INCLUDE_ROUTE_SERVICES", "true")
		})

		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
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
			os.Unsetenv("DEFAULT_TIMEOUT_IN_SECONDS")
			os.Unsetenv("CF_PUSH_TIMEOUT_IN_SECONDS")
			os.Unsetenv("LONG_CURL_TIMEOUT_IN_SECONDS")
			os.Unsetenv("BROKER_START_TIMEOUT_IN_SECONDS")
			os.Unsetenv("STATICFILE_BUILDPACK_NAME")
			os.Unsetenv("JAVA_BUILDPACK_NAME")
			os.Unsetenv("RUBY_BUILDPACK_NAME")
			os.Unsetenv("NODEJS_BUILDPACK_NAME")
			os.Unsetenv("GO_BUILDPACK_NAME")
			os.Unsetenv("PYTHON_BUILDPACK_NAME")
			os.Unsetenv("PHP_BUILDPACK_NAME")
			os.Unsetenv("BINARY_BUILDPACK_NAME")
			os.Unsetenv("INCLUDE_PRIVILEGED_CONTAINER_SUPPORT")
			os.Unsetenv("SKIP_SSO")

			os.Unsetenv("NODES")
			os.Unsetenv("INCLUDE_DIEGO_SSH")
			os.Unsetenv("INCLUDE_V3")
			os.Unsetenv("INCLUDE_DIEGO_DOCKER")
			os.Unsetenv("INCLUDE_BACKEND_COMPATIBILITY")
			os.Unsetenv("INCLUDE_SECURITY_GROUPS")
			os.Unsetenv("INCLUDE_LOGGING")
			os.Unsetenv("INCLUDE_OPERATOR")
			os.Unsetenv("INCLUDE_INTERNET_DEPENDENT")
			os.Unsetenv("INCLUDE_SERVICES")
			os.Unsetenv("INCLUDE_ROUTE_SERVICES")
		})
		It("Executes the command to run CATs", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Eventually(session.Out, 30).Should(gbytes.Say(
				"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes=5 -skipPackage=helpers -skip=NO_DIEGO_SUPPORT -keepGoing",
			))

			configBytes, err := ioutil.ReadFile(configJsonPath)
			Expect(err).NotTo(HaveOccurred())

			var config struct {
				Api                  string `json:"api"`
				AdminUser            string `json:"admin_user"`
				AdminPassword        string `json:"admin_password"`
				AppsDomain           string `json:"apps_domain"`
				SkipSSLValidation    bool   `json:"skip_ssl_validation"`
				UseHTTP              bool   `json:"use_http"`
				ExistingUser         string `json:"existing_user"`
				ExistingUserPassword string `json:"existing_user_password"`
				Backend              string `json:"backend"`

				PersistentAppHost      string `json:"persistent_app_host"`
				PersistentAppSpace     string `json:"persistent_app_space"`
				PersistentAppOrg       string `json:"persistent_app_org"`
				PersistentAppQuotaName string `json:"persistent_app_quota_name"`

				DefaultTimeout     int `json:"default_timeout"`
				CFPushTimeout      int `json:"cf_push_timeout"`
				LongCurlTimeout    int `json:"long_curl_timeout"`
				BrokerStartTimeout int `json:"broker_start_timeout"`

				StaticfileBuildpackName string `json:"staticfile_buildpack_name"`
				JavaBuildpackName       string `json:"java_buildpack_name"`
				RubyBuildpackName       string `json:"ruby_buildpack_name"`
				NodeJSBuilpackName      string `json:"nodejs_buildpack_name"`
				GoBuildpackName         string `json:"go_buildpack_name"`
				PythonBuildpackName     string `json:"python_buildpack_name"`
				PHPBuildpackName        string `json:"php_buildpack_name"`
				BinaryBuildpackName     string `json:"binary_buildpack_name"`

				IncludePrivilegedContainerSupport bool `json:"include_privileged_container_support"`
			}
			err = json.Unmarshal(configBytes, &config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Api).To(Equal("api.example.com"))
			Expect(config.AdminUser).To(Equal("admin-username"))
			Expect(config.AdminPassword).To(Equal("admin-password"))
			Expect(config.AppsDomain).To(Equal("apps.example.com"))
			Expect(config.SkipSSLValidation).To(BeTrue())
			Expect(config.UseHTTP).To(BeFalse())

			Expect(config.ExistingUser).To(Equal("existing-cats-user"))
			Expect(config.ExistingUserPassword).To(Equal("existing-cats-user-password"))
			Expect(config.Backend).To(Equal("diego"))
			Expect(config.PersistentAppHost).To(Equal("cats-app-host"))
			Expect(config.PersistentAppSpace).To(Equal("cats-app-space"))
			Expect(config.PersistentAppOrg).To(Equal("cats-app-org"))
			Expect(config.PersistentAppQuotaName).To(Equal("cats-app-quota"))

			Expect(config.DefaultTimeout).To(Equal(60))
			Expect(config.CFPushTimeout).To(Equal(120))
			Expect(config.LongCurlTimeout).To(Equal(180))
			Expect(config.BrokerStartTimeout).To(Equal(240))

			Expect(config.StaticfileBuildpackName).To(Equal("static-buildpack"))
			Expect(config.JavaBuildpackName).To(Equal("java-buildpack"))
			Expect(config.RubyBuildpackName).To(Equal("ruby-buildpack"))
			Expect(config.NodeJSBuilpackName).To(Equal("node-buildpack"))
			Expect(config.GoBuildpackName).To(Equal("go-buildpack"))
			Expect(config.PythonBuildpackName).To(Equal("python-buildpack"))
			Expect(config.PHPBuildpackName).To(Equal("php-buildpack"))
			Expect(config.BinaryBuildpackName).To(Equal("binary-buildpack"))

			Expect(config.IncludePrivilegedContainerSupport).To(Equal(true))
		})
	})
})
