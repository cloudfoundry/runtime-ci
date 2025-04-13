package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var pathToBinary string

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = SynchronizedBeforeSuite(
	func() []byte {
		path, err := gexec.Build("github.com/cloudfoundry/runtime-ci/util/update-manifest-releases")
		Expect(err).NotTo(HaveOccurred())

		return []byte(path)
	}, func(b []byte) {
		pathToBinary = string(b)
	},
)

var _ = SynchronizedAfterSuite(
	func() {},
	func() {
		gexec.CleanupBuildArtifacts()
	},
)
