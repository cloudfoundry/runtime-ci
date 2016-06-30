package commandgenerator_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/commandgenerator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commandgenerator", func() {
	var nodes int

	BeforeEach(func() {
		rand.Seed(time.Now().UTC().UnixNano())
		nodes = rand.Intn(100)
		os.Setenv("NODES", strconv.Itoa(nodes))
	})

	AfterEach(func() {
		os.Unsetenv("NODES")
	})

	It("Should generate a command to run CATS", func() {
		cmd := commandgenerator.GenerateCmd()
		Expect(cmd).To(Equal(fmt.Sprintf(
			"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes %d -skipPackage=helpers -skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing",
			nodes,
		)))
	})

	Context("when there are optional skipPackage env vars set", func() {
		BeforeEach(func() {
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

		It("should generate a command with the correct list of skipPackage flags", func() {
			cmd := commandgenerator.GenerateCmd()
			Expect(cmd).To(Equal(
				"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=helpers,ssh,v3,docker,backend_compatibility,security_groups,logging,operator,internet_dependent,services,route_services -skip=NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing",
			))
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
			cmd := commandgenerator.GenerateCmd()
			Expect(cmd).To(Equal(
				"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=helpers -skip=SSO|NO_DEA_SUPPORT|NO_DIEGO_SUPPORT -keepGoing",
			))
		})

		Context("when the backend is set to diego", func() {
			BeforeEach(func() {
				os.Setenv("BACKEND", "diego")
			})

			AfterEach(func() {
				os.Unsetenv("BACKEND")
			})

			It("should generate a command with the correct list of skip flags", func() {
				cmd := commandgenerator.GenerateCmd()
				Expect(cmd).To(Equal(
					"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=helpers -skip=SSO|NO_DEA_SUPPORT -keepGoing",
				))
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
				cmd := commandgenerator.GenerateCmd()
				Expect(cmd).To(Equal(
					"bin/test -r -slowSpecThreshold=120 -randomizeAllSpecs -nodes 5 -skipPackage=helpers -skip=SSO|NO_DIEGO_SUPPORT -keepGoing",
				))
			})
		})
	})
})
