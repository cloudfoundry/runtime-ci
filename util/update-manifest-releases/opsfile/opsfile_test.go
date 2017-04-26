package opsfile_test

import (
	"io/ioutil"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateReleases", func() {
	var (
		brokenBuildDir    string
		goodBuildDir      string
		noChangesBuildDir string

		originalOpsFile []byte
		desiredOpsFile  []byte
	)

	BeforeEach(func() {
		brokenBuildDir = "../fixtures/broken-build"
		goodBuildDir = "../fixtures/build"
		noChangesBuildDir = "../fixtures/nochanges-build"

		var err error

		desiredOpsFile, err = ioutil.ReadFile("../fixtures/desired_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())

		originalOpsFile, err = ioutil.ReadFile("../fixtures/original_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())

		manifest.ResetYAMLMarshal()
	})

	It("updates only releases with different shas without modifying the rest of the file", func() {
		releaseNames := []string{"release1", "release2"}

		updatedOpsFile, changes, err := opsfile.UpdateReleases(releaseNames, goodBuildDir, originalOpsFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(updatedOpsFile)).To(Equal(string(desiredOpsFile)))
		Expect(changes).To(Equal("Updated release2-release"))
	})

	It("provides a default commit message if no version updates were performed", func() {
		releaseNames := []string{"release1", "release2"}

		_, changes, err := opsfile.UpdateReleases(releaseNames, noChangesBuildDir, originalOpsFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(changes).To(Equal("No release updates"))
	})

	Context("failure cases", func() {
		It("returns errors instead of panicking when url is missing", func() {
			releases := []string{"missing-url"}

			_, _, err := opsfile.UpdateReleases(releases, brokenBuildDir, originalOpsFile)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-url-release/url: no such file or directory"))
		})

		It("returns errors instead of panicking when version is missing", func() {
			releases := []string{"missing-version"}

			_, _, err := opsfile.UpdateReleases(releases, brokenBuildDir, originalOpsFile)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-version-release/version: no such file or directory"))
		})

		It("returns errors instead of panicking when sha1 is missing", func() {
			releases := []string{"missing-sha1"}

			_, _, err := opsfile.UpdateReleases(releases, brokenBuildDir, originalOpsFile)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-sha1-release/sha1: no such file or directory"))
		})

		It("returns an error when the manifest is not valid yaml", func() {
			releases := []string{"release1", "release2"}

			originalOpsFile := []byte(`
%%%
releases:
%%%
`)
			_, _, err := opsfile.UpdateReleases(releases, goodBuildDir, originalOpsFile)
			Expect(err).To(MatchError(ContainSubstring("could not find expected directive name")))
		})
	})
})
