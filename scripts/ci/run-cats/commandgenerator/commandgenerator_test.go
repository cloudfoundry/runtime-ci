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
			BeforeEach(func() {
				env.GetBooleanReturns(true, nil)
			})

			It("should generate a command with the correct list of skipPackage flags", func() {
				cmd, args, err := commandgenerator.GenerateCmd(env)
				Expect(err).NotTo(HaveOccurred())
				Expect(cmd).To(Equal(
					"bin/test",
				))

				Expect(strings.Join(args, " ")).To(Equal(
					fmt.Sprintf("-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes %d -skipPackage=helpers -skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing", nodes)))
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
			BeforeEach(func() {
				os.Setenv("NODES", "5")
				os.Setenv("SKIP_SSO", "true")
			})

			AfterEach(func() {
				os.Unsetenv("NODES")
				os.Unsetenv("SKIP_SSO")
			})

			It("should generate a command with the correct list of skip flags", func() {
				cmd, args, err := commandgenerator.GenerateCmd(env)
				Expect(err).NotTo(HaveOccurred())
				Expect(cmd).To(Equal(
					"bin/test",
				))

				Expect(strings.Join(args, " ")).To(Equal(
					"-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 -skip=SSO|NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing"))
			})

			Context("when the backend is set to diego", func() {
				BeforeEach(func() {
					os.Setenv("BACKEND", "diego")
				})

				AfterEach(func() {
					os.Unsetenv("BACKEND")
				})

				It("should generate a command with the correct list of skip flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))

					Expect(strings.Join(args, " ")).To(Equal(
						"-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 -skip=SSO|NO_DIEGO_SUPPORT -keepGoing"))
				})

			})

			Context("when the backend is set to dea", func() {
				BeforeEach(func() {
					os.Setenv("BACKEND", "dea")
				})

				AfterEach(func() {
					os.Unsetenv("BACKEND")
				})

				It("should generate a command with the correct list of skip flags", func() {
					cmd, args, err := commandgenerator.GenerateCmd(env)
					Expect(err).NotTo(HaveOccurred())
					Expect(cmd).To(Equal(
						"bin/test",
					))
					Expect(strings.Join(args, " ")).To(Equal(
						"-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=backend_compatibility,docker,helpers,internet_dependent,logging,operator,route_services,security_groups,services,ssh,v3 -skip=SSO|NO_DEA_SUPPORT -keepGoing"))
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
