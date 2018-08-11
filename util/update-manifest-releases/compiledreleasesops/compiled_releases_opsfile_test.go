package compiledreleasesops_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/compiledreleasesops"
	"io/ioutil"
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

		updatedOpsFile, commitMessage, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)

		Expect(err).ToNot(HaveOccurred())
		Expect(commitMessage).To(Equal("Updated compiled releases with test 0.1.0"))
		Expect(updatedOpsFile).To(MatchYAML(desiredOpsFile))
	})

	It("returns error when there's more than one compiled release tarball", func() {
		releaseNames := []string{"more-than-1"}

		_, _, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("expected to find exactly 1 compiled release tarball"))
	})

	It("returns error when release folder doesn't have a version file", func() {
		releaseNames := []string{"no-version"}

		_, _, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("could not find necessary release info:"))
	})

	It("returns error when it could not find desired release in the ops file", func() {
		releaseNames := []string{"extraneous"}

		_, _, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("could not find release 'extraneous' in the ops file"))
	})
})