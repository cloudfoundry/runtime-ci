package commandgenerator_test

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator/commandgeneratorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commandgenerator", func() {
	var nodes int
	var env *commandgeneratorfakes.FakeEnvironment

	BeforeEach(func() {
		env = &commandgeneratorfakes.FakeEnvironment{}
		nodes = 10
		env.GetIntegerReturnsFor("NODES", nodes, nil)

	})

	Context("When the path to CATs is set", func() {
		BeforeEach(func() {
			env.GetStringReturnsFor("CATS_PATH", ".")
		})

		It("Should generate a command to run CATS", func() {
			cmd, args, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("bin/test"))

			Expect(strings.Join(args, " ")).To(Equal(
				fmt.Sprintf("-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes=%d -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 -skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing", nodes),
			))

			env.GetStringReturnsFor("CATS_PATH", "/path/to/cats")
			cmd, _, err = commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/path/to/cats/bin/test"))
		})

		Context("when the env returns an error fetching the number of nodes", func() {
			var expectedError error
			BeforeEach(func() {
				expectedError = fmt.Errorf("some error")
				env.GetIntegerReturnsFor("NODES", 0, expectedError)
			})

			It("propogates the error", func() {
				_, _, err := commandgenerator.GenerateCmd(env)
				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when the node count is unset", func() {
			BeforeEach(func() {
				env.GetIntegerReturnsFor("NODES", 0, nil)
			})
			It("sets the default node count", func() {
				_, args, _ := commandgenerator.GenerateCmd(env)
				Expect(args).To(ContainElement("-nodes=2"))
			})
		})

		Context("when there are optional skipPackage env vars set", func() {
			Context("to true", func() {
				BeforeEach(func() {
					env.GetBooleanReturnsFor("INCLUDE_SSO", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_BACKEND_COMPATIBILITY", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_DIEGO_DOCKER", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_INTERNET_DEPENDENT", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_LOGGING", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_OPERATOR", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_ROUTE_SERVICES", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_SECURITY_GROUPS", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_SERVICES", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_DIEGO_SSH", true, nil)
					env.GetBooleanReturnsFor("INCLUDE_V3", true, nil)
				})

				It("should generate a command with the correct list of skipPackage flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(strings.Join(args, " ")).To(ContainSubstring(" -skipPackage=helpers "))
				})
			})

			Context("to false", func() {
				BeforeEach(func() {
					env.GetBooleanReturnsFor("INCLUDE_SSO", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_BACKEND_COMPATIBILITY", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_DIEGO_DOCKER", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_INTERNET_DEPENDENT", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_LOGGING", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_OPERATOR", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_ROUTE_SERVICES", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_SECURITY_GROUPS", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_SERVICES", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_DIEGO_SSH", false, nil)
					env.GetBooleanReturnsFor("INCLUDE_V3", false, nil)
				})

				It("should generate a command with the correct list of skipPackage flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(strings.Join(args, " ")).To(ContainSubstring(
						" -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 ",
					))
				})
			})

			Context("when the env returns an error", func() {
				var expectedError error
				BeforeEach(func() {
					expectedError = fmt.Errorf("some error")
					env.GetBooleanReturnsFor("INCLUDE_V3", false, expectedError)
				})

				It("propogates the error", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err).To(Equal(expectedError))
				})
			})
		})

		Context("when there are optional skip env vars set", func() {
			Context("to true", func() {
				BeforeEach(func() {
					env.GetBooleanReturnsFor("SKIP_SSO", true, nil)
				})

				It("generates a command that skips tests with the given tag", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).To(ContainElement(ContainSubstring("-skip=SSO|")))
				})
			})

			Context("to false", func() {
				BeforeEach(func() {
					env.GetBooleanReturnsFor("SKIP_SSO", false, nil)
				})

				It("generates a command that does not include the given tag in the skips", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).ToNot(ContainElement(ContainSubstring("-skip=SSO|")))
				})
			})

			Context("and the env returns an error", func() {
				var expectedError error
				BeforeEach(func() {
					expectedError = fmt.Errorf("some error")
					env.GetBooleanReturnsFor("SKIP_SSO", false, expectedError)
				})

				It("propogates the error", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err).To(Equal(expectedError))
					Expect(env.GetBooleanCallCountFor("SKIP_SSO"))
				})
			})

			Context("when the backend is set to diego", func() {
				BeforeEach(func() {
					env.GetStringReturnsFor("BACKEND", "diego")
				})

				It("should generate a command that skips NO_DIEGO_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).To(ContainElement("-skip=NO_DIEGO_SUPPORT"))
				})

			})

			Context("when the backend is set to dea", func() {
				BeforeEach(func() {
					env.GetStringReturnsFor("BACKEND", "dea")
				})

				It("should generate a command that skips NO_DEA_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))
					Expect(args).To(ContainElement("-skip=NO_DEA_SUPPORT"))
				})
			})

			Context("when the backend isn't set", func() {
				It("should generate a command that skips bosh NO_DIEGO_SUPPORT and NO_DEA_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))
					Expect(args).To(ContainElement("-skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT"))
				})
			})

			Context("when the backend is not a valid value", func() {
				BeforeEach(func() {
					env.GetStringReturnsFor("BACKEND", "bogus")
				})
				It("should generate a command that skips bosh NO_DIEGO_SUPPORT and NO_DEA_SUPPORT", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid environment variable: 'BACKEND' was 'bogus', but must be 'diego', 'dea', or empty"))
				})
			})
		})
	})

	Context("When the path to CATS isn't explicitly provided", func() {
		BeforeEach(func() {
			env.GetStringReturnsFor("GOPATH", "/go")
		})

		It("Should return a sane default command path for use in Concourse", func() {
			cmd, _, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/go/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"))
		})
	})
})
