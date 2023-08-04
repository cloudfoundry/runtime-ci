package compiledreleasesops_test

import (
	"os"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/compiledreleasesops"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("UpdateCompiledReleases", func() {
	var (
		compiledReleaseBuildDir string

		originalOpsFile []byte
		desiredOpsFile  []byte

		err error
	)

	BeforeEach(func() {
		compiledReleaseBuildDir = "../fixtures/build-with-compiled-release"

		desiredOpsFile, err = os.ReadFile("../fixtures/updated_compiled_releases_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())

		originalOpsFile, err = os.ReadFile("../fixtures/original_compiled_releases_ops_file.yml")
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

	It("adds release if it cannot be found in the ops file", func() {
		releaseNames := []string{"extraneous"}
		desiredOpsFile, err = os.ReadFile("../fixtures/updated_compiled_releases_ops_file_with_new_release.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedOpsFile, commitMessage, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())
		Expect(commitMessage).To(Equal("Updated compiled releases with extraneous 0.1.0"))
		Expect(updatedOpsFile).To(MatchYAML(desiredOpsFile))
	})

	It("adds stemcell section if the release doesn't have one in the existing compiled releases ops file", func() {
		releaseNames := []string{"no-stemcell-section"}
		originalOpsFile := `- type: replace
  path: /releases/name=no-stemcell-section
  value:
    name: no-stemcell-section
    sha1: 5ee0dfe1f1b9acd14c18863061268f4156c291a4
    url: https://storage.googleapis.com/cf-deployment-compiled-releases/no-stemcell-section-0.0.0-cute-stemcell-0.0-20180808-195254-497840039.tgz
    version: 0.0.1
`
		desiredOpsFile := `- type: replace
  path: /releases/name=no-stemcell-section
  value:
    name: no-stemcell-section
    sha1: 02573f83a7f467e55a7bb49424e80f541288a041
    stemcell:
      os: awesome-stemcell
      version: "1.0"
    url: https://storage.googleapis.com/cf-deployment-compiled-releases/no-stemcell-section-0.3.0-awesome-stemcell-1.0-20180808-195254-497840039.tgz
    version: 0.3.0
`
		updatedOpsFile, commitMessage, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, []byte(originalOpsFile), yaml.Marshal, yaml.Unmarshal)

		Expect(err).NotTo(HaveOccurred())
		Expect(commitMessage).To(Equal("Updated compiled releases with no-stemcell-section 0.3.0"))
		Expect(string(updatedOpsFile)).To(Equal(desiredOpsFile))
	})
})
