package commandgenerator_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator/commandgeneratorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commandgenerator", func() {
	var nodes int
	var env *commandgeneratorfakes.FakeEnvironment

	BeforeEach(func() {
		rand.Seed(time.Now().UTC().UnixNano())
		nodes = rand.Intn(100)
		os.Setenv("NODES", strconv.Itoa(nodes))

		env = &commandgeneratorfakes.FakeEnvironment{}
	})

	AfterEach(func() {
		os.Unsetenv("NODES")
	})

	Context("When the path to CATs is set", func() {
		BeforeEach(func() {
			os.Setenv("CATS_PATH", ".")
		})

		AfterEach(func() {
			os.Unsetenv("CATS_PATH")
		})

		It("Should generate a command to run CATS", func() {
			cmd, args, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("bin/test"))

			Expect(strings.Join(args, " ")).To(Equal(
				fmt.Sprintf("-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes %d -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 -skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing", nodes),
			))

			os.Setenv("CATS_PATH", "/path/to/cats")
			cmd, _, err = commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/path/to/cats/bin/test"))
		})

		Context("when there are optional skipPackage env vars set", func() {
			Context("to true", func() {
				BeforeEach(func() {
					env.GetBooleanReturns(true, nil)
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
					env.GetBooleanReturns(false, nil)
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
					env.GetBooleanReturns(false, expectedError)
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
					env.GetBooleanReturns(true, nil)
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
					env.GetBooleanReturns(false, nil)
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
					env.GetBooleanStub = func(varName string) (bool, error) {
						switch varName {
						case "SKIP_SSO":
							return false, expectedError
						default:
							return false, nil
						}
					}
				})

				It("propogates the error", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err).To(Equal(expectedError))

					found := false
					for i := 0; i < env.GetBooleanCallCount(); i++ {
						varName := env.GetBooleanArgsForCall(i)
						if varName == "SKIP_SSO" {
							found = true
						}
					}
					Expect(found).To(BeTrue())
				})
			})

			Context("when the backend is set to diego", func() {
				BeforeEach(func() {
					env.GetStringReturns("diego")
				})

				It("should generate a command that skips NO_DIEGO_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(args).To(ContainElement("-skip=NO_DIEGO_SUPPORT"))
					Expect(env.GetStringArgsForCall(0)).To(Equal("BACKEND"))
				})

			})

			Context("when the backend is set to dea", func() {
				BeforeEach(func() {
					env.GetStringReturns("dea")
				})

				It("should generate a command that skips NO_DEA_SUPPORT", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))
					Expect(args).To(ContainElement("-skip=NO_DEA_SUPPORT"))
					Expect(env.GetStringArgsForCall(0)).To(Equal("BACKEND"))
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
					Expect(env.GetStringArgsForCall(0)).To(Equal("BACKEND"))
				})
			})

			Context("when the backend is not a valid value", func() {
				BeforeEach(func() {
					env.GetStringReturns("bogus")
				})
				It("should generate a command that skips bosh NO_DIEGO_SUPPORT and NO_DEA_SUPPORT", func() {
					_, _, err := commandgenerator.GenerateCmd(env)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid environment variable: 'BACKEND' was 'bogus', but must be 'diego', 'dea', or empty"))
					Expect(env.GetStringArgsForCall(0)).To(Equal("BACKEND"))
				})
			})
		})
	})

	Context("When the path to CATS isn't explicitly provided", func() {
		BeforeEach(func() {
			os.Setenv("GOPATH", "/go")
		})

		AfterEach(func() {
			os.Unsetenv("GOPATH")
		})

		It("Should return a sane default command path for use in Concourse", func() {
			cmd, _, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/go/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"))
		})
	})
})
