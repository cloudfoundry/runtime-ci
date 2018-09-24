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

	It("adds release if it cannot be found in the ops file", func() {
		releaseNames := []string{"extraneous"}
		desiredOpsFile, err = ioutil.ReadFile("../fixtures/updated_compiled_releases_ops_file_with_new_release.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedOpsFile, commitMessage, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())
		Expect(commitMessage).To(Equal("Updated compiled releases with extraneous 0.1.0"))
		Expect(updatedOpsFile).To(MatchYAML(desiredOpsFile))
	})

	It("adds stemcell section if the release doesn't have one in the existing compiled releases ops file", func() {
		releaseNames := []string{"no-stemcell-section"}
		originalOpsFile := `
- path: /releases/name=no-stemcell-section/url
  type: replace
  value: https://storage.googleapis.com/cf-deployment-compiled-releases/no-stemcell-section-0.0.0-cute-stemcell-0.0-20180808-195254-497840039.tgz
- path: /releases/name=no-stemcell-section/version
  type: replace
  value: 0.0.1
- path: /releases/name=no-stemcell-section/sha1
  type: replace
  value: 5ee0dfe1f1b9acd14c18863061268f4156c291a4
`
		desiredOpsFile := `
- path: /releases/name=no-stemcell-section/url
  type: replace
  value: https://storage.googleapis.com/cf-deployment-compiled-releases/no-stemcell-section-0.3.0-awesome-stemcell-1.0-20180808-195254-497840039.tgz
- path: /releases/name=no-stemcell-section/version
  type: replace
  value: 0.3.0
- path: /releases/name=no-stemcell-section/sha1
  type: replace
  value: 02573f83a7f467e55a7bb49424e80f541288a041
- path: /releases/name=no-stemcell-section/stemcell?
  type: replace
  value:
    os: awesome-stemcell
    version: "1.0"
`
		updatedOpsFile, commitMessage, err := compiledreleasesops.UpdateCompiledReleases(releaseNames, compiledReleaseBuildDir, []byte(originalOpsFile), yaml.Marshal, yaml.Unmarshal)

		Expect(err).NotTo(HaveOccurred())
		Expect(commitMessage).To(Equal("Updated compiled releases with no-stemcell-section 0.3.0"))
		Expect(updatedOpsFile).To(MatchYAML(desiredOpsFile))
	})
})