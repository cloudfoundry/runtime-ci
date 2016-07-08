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
	AfterEach(func() {
		os.Remove(os.Getenv("PWD") + "/integration_config.json")
	})

	Context("when no envvars are set", func() {

		It("Exits 1 and prints an error regarding missing 'required' env vars", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(1))
			Eventually(session.Out, 30).Should(gbytes.Say(`Missing required environment variables:
CF_API
CF_ADMIN_USER
CF_ADMIN_PASSWORD
CF_APPS_DOMAIN`,
			))

			Expect(string(session.Out.Contents())).NotTo(ContainSubstring("CONFIG="))
			Expect(os.Getenv("PWD") + "/integration_config.json").NotTo(BeARegularFile())
		})
	})

	Context("when some but not all required envvars are set", func() {
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
			Eventually(session.Out, 30).Should(gbytes.Say(`Missing required environment variables:
CF_ADMIN_USER
CF_ADMIN_PASSWORD
CF_APPS_DOMAIN`,
			))
			Expect(os.Getenv("PWD") + "/integration_config.json").NotTo(BeARegularFile())

		})
	})

	Context("When invalid TIMEOUT env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("CATS_PATH", ".")
			os.Setenv("CF_API", "non-empty-value")
			os.Setenv("CF_ADMIN_USER", "non-empty-value")
			os.Setenv("CF_ADMIN_PASSWORD", "non-empty-value")
			os.Setenv("CF_APPS_DOMAIN", "non-empty-value")
			os.Setenv("SKIP_SSL_VALIDATION", "true")
			os.Setenv("USE_HTTP", "true")

			os.Setenv("DEFAULT_TIMEOUT_IN_SECONDS", "not a timeout")
		})
		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
			os.Unsetenv("SKIP_SSL_VALIDATION")
			os.Unsetenv("USE_HTTP")

			os.Unsetenv("DEFAULT_TIMEOUT_IN_SECONDS")
		})

		It("Should exit with an appropriate error message and exit code", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(1))
			Eventually(session.Err, 30).Should(gbytes.Say(`Invalid env var 'DEFAULT_TIMEOUT_IN_SECONDS' only allows positive integers`))
			Expect(os.Getenv("PWD") + "/integration_config.json").NotTo(BeARegularFile())
		})
	})

	Context("When invalid boolean env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("CATS_PATH", ".")
			os.Setenv("CF_API", "non-empty-value")
			os.Setenv("CF_ADMIN_USER", "non-empty-value")
			os.Setenv("CF_ADMIN_PASSWORD", "non-empty-value")
			os.Setenv("CF_APPS_DOMAIN", "non-empty-value")
			os.Setenv("SKIP_SSL_VALIDATION", "true")
			os.Setenv("USE_HTTP", "this is not a boolean this is only a tribute")
		})

		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
			os.Unsetenv("SKIP_SSL_VALIDATION")
			os.Unsetenv("USE_HTTP")
		})

		It("Should exit with an appropriate error message and exit code", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(1))
			Eventually(session.Err, 30).Should(gbytes.Say(`Invalid env var 'USE_HTTP' only accepts booleans`))
			Expect(os.Getenv("PWD") + "/integration_config.json").NotTo(BeARegularFile())
		})
	})

	Context("When all required env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("CATS_PATH", ".")
			os.Setenv("CF_API", "non-empty-value")
			os.Setenv("CF_ADMIN_USER", "non-empty-value")
			os.Setenv("CF_ADMIN_PASSWORD", "non-empty-value")
			os.Setenv("CF_APPS_DOMAIN", "non-empty-value")
			os.Setenv("SKIP_SSL_VALIDATION", "true")
			os.Setenv("USE_HTTP", "true")
		})
		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
			os.Unsetenv("CF_API")
			os.Unsetenv("CF_ADMIN_USER")
			os.Unsetenv("CF_ADMIN_PASSWORD")
			os.Unsetenv("CF_APPS_DOMAIN")
			os.Unsetenv("SKIP_SSL_VALIDATION")
			os.Unsetenv("USE_HTTP")
		})

		It("Sets CONFIG envvar to the path for the generated config file", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))

			Expect(string(session.Out.Contents())).To(ContainSubstring("CONFIG=" + os.Getenv("PWD") + "/integration_config.json"))
		})

		It("Writes a config file for CATs to use", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Expect(os.Getenv("PWD") + "/integration_config.json").To(BeARegularFile())
		})

		It("Executes the command to run CATs", func() {
			command := exec.Command(binPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Eventually(session.Out, 30).Should(gbytes.Say(
				"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 4 -skipPackage=helpers -skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing",
			))
		})
	})

	Context("When all supported env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("CATS_PATH", ".")
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
			os.Setenv("BINARY_BUILDPACK_NAME", "binary-builpack")

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
				"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=helpers -skip=NO_DIEGO_SUPPORT -keepGoing",
			))
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
})
