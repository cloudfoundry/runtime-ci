package opsfile_test

import (
	"errors"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/opsfile"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateReleases", func() {
	var (
		brokenBuildDir    string
		goodBuildDir      string
		noChangesBuildDir string

		originalOpsFile []byte
	)

	BeforeEach(func() {
		brokenBuildDir = "../fixtures/broken-build"
		goodBuildDir = "../fixtures/build"
		noChangesBuildDir = "../fixtures/nochanges-build"

		var err error

		originalOpsFile, err = os.ReadFile("../fixtures/original_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())
	})

	It("does not update release when opsfile is removing that release", func() {
		originalOpsFile, err := os.ReadFile("../fixtures/original_release_removal_opsfile.yml")
		Expect(err).NotTo(HaveOccurred())

		desiredOpsFile, err := os.ReadFile("../fixtures/updated_release_removal_opsfile.yml")
		Expect(err).NotTo(HaveOccurred())

		releaseNames := []string{"release1"}
		updatedOpsFile, changes, err := opsfile.UpdateReleases(releaseNames, "../fixtures/build", originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(updatedOpsFile).To(MatchYAML(desiredOpsFile))
		Expect(changes).To(Equal("No opsfile release updates"))
	})

	It("updates releases when opsfile does not use append syntax", func() {
		releaseNames := []string{"non-append"}

		originalOpsFile, err := os.ReadFile("../fixtures/original_non_append_opsfile.yml")
		Expect(err).NotTo(HaveOccurred())

		desiredOpsFile, err := os.ReadFile("../fixtures/updated_non_append_opsfile.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedOpsFile, changes, err := opsfile.UpdateReleases(releaseNames, "../fixtures/build", originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(updatedOpsFile).To(MatchYAML(desiredOpsFile))
		Expect(changes).To(Equal("Updated ops file(s) with non-append-release updated-non-append-version"))
	})

	It("updates releases with different shas", func() {
		releaseNames := []string{"release1", "release2"}

		desiredOpsFile, err := os.ReadFile("../fixtures/updated_sha_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedOpsFile, changes, err := opsfile.UpdateReleases(releaseNames, "../fixtures/build-with-updated-sha", originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(updatedOpsFile)).To(MatchYAML(desiredOpsFile))
		Expect(changes).To(Equal("Updated ops file(s) with release2-release original-release2-version"))
	})

	It("updates releases with different versions", func() {
		releaseNames := []string{"release1", "release2"}

		desiredOpsFile, err := os.ReadFile("../fixtures/updated_version_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedOpsFile, changes, err := opsfile.UpdateReleases(releaseNames, "../fixtures/build-with-updated-version", originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(updatedOpsFile)).To(MatchYAML(string(desiredOpsFile)))
		Expect(changes).To(Equal("Updated ops file(s) with release2-release updated-release2-version"))
	})

	It("updates releases with different urls", func() {
		releaseNames := []string{"release1", "release2"}

		desiredOpsFile, err := os.ReadFile("../fixtures/updated_url_ops_file.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedOpsFile, changes, err := opsfile.UpdateReleases(releaseNames, "../fixtures/build-with-updated-url", originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(updatedOpsFile)).To(MatchYAML(desiredOpsFile))
		Expect(changes).To(Equal("Updated ops file(s) with release2-release original-release2-version"))
	})

	It("provides a default commit message if no version updates were performed", func() {
		releaseNames := []string{"release1", "release2"}

		_, changes, err := opsfile.UpdateReleases(releaseNames, noChangesBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(changes).To(Equal("No opsfile release updates"))
	})

	Context("failure cases", func() {
		It("returns errors instead of panicking when url is missing", func() {
			releases := []string{"missing-url"}

			_, _, err := opsfile.UpdateReleases(releases, brokenBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-url-release/url: no such file or directory"))
		})

		It("returns errors instead of panicking when version is missing", func() {
			releases := []string{"missing-version"}

			_, _, err := opsfile.UpdateReleases(releases, brokenBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-version-release/version: no such file or directory"))
		})

		It("returns errors instead of panicking when sha1 is missing", func() {
			releases := []string{"missing-sha1"}

			_, _, err := opsfile.UpdateReleases(releases, brokenBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-sha1-release/sha1: no such file or directory"))
		})

		It("returns an error when the manifest is not valid yaml", func() {
			releases := []string{"release1", "release2"}

			originalOpsFile := []byte(`
%%%
releases:
%%%
`)
			_, _, err := opsfile.UpdateReleases(releases, goodBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(ContainSubstring("could not find expected directive name")))
		})

		It("does not add a `value: null` field to remove operations", func() {
			releases := []string{"release1"}

			originalOpsFile := []byte(`
- type: remove
  path: /stemcell

- type: replace
  path: /releases/-
  value:
    name: release1
    version: foo
`)
			updatedOpsFile, _, err := opsfile.UpdateReleases(releases, goodBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)

			Expect(err).ToNot(HaveOccurred())
			Expect(updatedOpsFile).ToNot(ContainSubstring("null"))
		})

		It("returns an error when the yaml marshaller fails", func() {
			failingMarshalFunc := func(interface{}) ([]byte, error) {
				return nil, errors.New("failed to marshal yaml")
			}
			releases := []string{"release1", "release2"}

			_, _, err := opsfile.UpdateReleases(releases, goodBuildDir, originalOpsFile, failingMarshalFunc, yaml.Unmarshal)
			Expect(err).To(MatchError("failed to marshal yaml"))
		})

		It("returns an error when the yaml unmarshaller fails", func() {
			failingUnmarshalFunc := func([]byte, interface{}) error {
				return errors.New("failed to unmarshal yaml")
			}
			releases := []string{"release1", "release2"}

			_, _, err := opsfile.UpdateReleases(releases, goodBuildDir, originalOpsFile, yaml.Marshal, failingUnmarshalFunc)
			Expect(err).To(MatchError("failed to unmarshal yaml"))
		})

		It("returns an error when the original ops file does not contain an expected release", func() {
			releases := []string{"fun-times"}
			originalOpsFile := []byte(`
- type: replace
  path: /releases/-
  value:
    name: sad-times
    version: 1.0.0
`)
			_, _, err := opsfile.UpdateReleases(releases, goodBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError("opsfile does not contain release named fun-times"))
		})

		It("returns an error when the release name array is nil or empty", func() {
			_, _, err := opsfile.UpdateReleases(nil, goodBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError("releaseNames provided to UpdateReleases must contain at least one release name"))

			_, _, err = opsfile.UpdateReleases([]string{}, goodBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError("releaseNames provided to UpdateReleases must contain at least one release name"))
		})

		It("returns an error when the ops file has separate ops for each release value", func() {
			releases := []string{"fun-times"}
			originalOpsFile := []byte(`
- path: /releases/name=test/url
  type: replace
  value: release-url
- path: /releases/name=test/version
  type: replace
  value: 0.0.0
- path: /releases/name=test/sha1
  type: replace
  value: sha256:4ee0dfe1f1b9acd14c18863061268f4156c291a4
`)
			_, _, err := opsfile.UpdateReleases(releases, goodBuildDir, originalOpsFile, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(opsfile.BadReleaseOpsFormatErrorMessage))
		})
	})
})
