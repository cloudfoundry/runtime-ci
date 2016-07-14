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
})
