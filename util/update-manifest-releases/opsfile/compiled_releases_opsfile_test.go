package opsfile_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
	"gopkg.in/yaml.v2"
)

var _ = Describe("UpdateCompiledReleases", func() {
	var (
		compiledReleaseBuildDir    string

		originalOpsFile []byte
		desiredOpsFile  []byte

		err error
	)

	BeforeEach(func() {
		compiledReleaseBuildDir = "../fixtures/build-with-compiled-release"

		desiredOpsFile, err = ioutil.ReadFile("../fixtures/updated_compiled_releases_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())

		originalOpsFile, err = ioutil.ReadFile("../fixtures/original_compiled_releases_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())
	})

	It("updates compiled releases ops file for the desired release", func() {
		releaseNames := []string{"test"}

		updatedOpsFile, commitMessage, err := opsfile.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)

		Expect(err).ToNot(HaveOccurred())
		Expect(commitMessage).To(BeEmpty())
		Expect(updatedOpsFile).To(MatchYAML(desiredOpsFile))
	})
})