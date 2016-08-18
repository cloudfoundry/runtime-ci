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
		env.GetNodesReturns(nodes, nil)
	})

	Context("When the path to CATs is set", func() {
		BeforeEach(func() {
			env.GetCatsPathReturns(".")
		})

		It("Should generate a command to run CATS", func() {
			cmd, args, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("bin/test"))

			Expect(strings.Join(args, " ")).To(Equal(
				fmt.Sprintf("-r -slowSpecThreshold=120 -randomizeAllSpecs -nodes=%d -skipPackage=helpers -keepGoing", nodes),
			))

			env.GetCatsPathReturns("/path/to/cats")
			cmd, _, err = commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/path/to/cats/bin/test"))
		})

		Context("when the node count is unset", func() {
			BeforeEach(func() {
				env.GetNodesReturns(0, nil)
			})
			It("sets the default node count", func() {
				_, args, _ := commandgenerator.GenerateCmd(env)
				Expect(args).To(ContainElement("-nodes=2"))
			})
		})
	})

	Context("When the path to CATS isn't explicitly provided", func() {
		BeforeEach(func() {
			env.GetGoPathReturns("/go")
		})

		It("Should return a sane default command path for use in Concourse", func() {
			cmd, _, err := commandgenerator.GenerateCmd(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmd).To(Equal("/go/src/github.com/cloudfoundry/cf-acceptance-tests/bin/test"))
		})
	})
})
