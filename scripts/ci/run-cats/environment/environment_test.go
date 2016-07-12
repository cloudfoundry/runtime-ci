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

			It("returns an error", func() {
				env := environment.New()
				_, err := env.GetBoolean("MY_ENV_VAR")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid environment variable: 'MY_ENV_VAR' must be a boolean 'true' or 'false'"))
			})
		})
	})
})
