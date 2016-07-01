package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var binPath string

var _ = BeforeSuite(func() {
	By("Compiling binary")
	var err error
	binPath, err = gexec.Build("github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestRunCats(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RunCats Suite")
}
