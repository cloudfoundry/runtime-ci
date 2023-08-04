package manifest_test

import (
	"errors"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateReleases", func() {
	var (
		brokenBuildDir    string
		goodBuildDir      string
		noChangesBuildDir string

		cfDeploymentManifest []byte
	)

	BeforeEach(func() {
		brokenBuildDir = "../fixtures/broken-build"
		goodBuildDir = "../fixtures/build"
		noChangesBuildDir = "../fixtures/nochanges-build"

		var err error

		cfDeploymentManifest, err = os.ReadFile("../fixtures/cf-deployment.yml")
		Expect(err).NotTo(HaveOccurred())
	})

	It("updates the releases without modifying the rest and returns the list of changes when the sha changes", func() {
		releases := []string{"release1", "release2"}
		updatedReleasesFixture, err := os.ReadFile("../fixtures/updated_sha_releases.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateReleases(releases, "../fixtures/build-with-updated-sha", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleases := updatedManifest[updatedManifestReleasesIndex:]

		Expect(string(cfDeploymentPreamble)).To(Equal(string(updatedManifestPreamble)), "the preamble was changed by running the program")
		Expect(string(updatedManifestReleases)).To(Equal(string(updatedReleasesFixture)))

		Expect(changes).To(Equal("Updated manifest with release2-release original-release2-version"))
	})

	It("updates the releases without modifying the rest and returns the list of changes when the version changes", func() {
		releases := []string{"release1", "release2"}
		updatedReleasesFixture, err := os.ReadFile("../fixtures/updated_version_releases.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateReleases(releases, "../fixtures/build-with-updated-version", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleases := updatedManifest[updatedManifestReleasesIndex:]

		Expect(string(cfDeploymentPreamble)).To(Equal(string(updatedManifestPreamble)), "the preamble was changed by running the program")
		Expect(string(updatedManifestReleases)).To(Equal(string(updatedReleasesFixture)))

		Expect(changes).To(Equal("Updated manifest with release2-release updated-release2-version"))
	})

	It("updates the releases without modifying the rest and returns the list of changes when the url changes", func() {
		releases := []string{"release1", "release2"}
		updatedReleasesFixture, err := os.ReadFile("../fixtures/updated_url_releases.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateReleases(releases, "../fixtures/build-with-updated-url", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleases := updatedManifest[updatedManifestReleasesIndex:]

		Expect(string(cfDeploymentPreamble)).To(Equal(string(updatedManifestPreamble)), "the preamble was changed by running the program")
		Expect(string(updatedManifestReleases)).To(Equal(string(updatedReleasesFixture)))

		Expect(changes).To(Equal("Updated manifest with release2-release original-release2-version"))
	})

	It("provides a default commit message if no version updates were performed", func() {
		releases := []string{"release1", "release2"}
		_, changes, err := manifest.UpdateReleases(releases, noChangesBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(changes).To(Equal("No manifest release or stemcell version updates"))
	})

	It("adds them to resulting list of releases when there are updates to the releases that are not in the manifest releases", func() {
		updateReleases := []string{"release1"}
		cfDeploymentManifest := []byte(`
name: my-deployment
releases:
  - name: fooRelease
stemcells:
`)
		resultingManifest, _, err := manifest.UpdateReleases(updateReleases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).ToNot(HaveOccurred())

		var releases manifest.Manifest
		err = yaml.Unmarshal(resultingManifest, &releases)
		Expect(err).ToNot(HaveOccurred())

		Expect(releases.Releases).To(ContainElement(common.Release{
			Name:    "release1",
			URL:     "original-release1-url",
			Version: "original-release1-version",
			SHA1:    "original-release1-sha1",
		}))

	})

	It("adds all the releases to the commit message if no releases exist", func() {
		releases := []string{"release1", "release2"}
		cfDeploymentManifest := []byte(`
name: my-deployment
releases:
stemcells:
`)

		_, changes, err := manifest.UpdateReleases(releases, noChangesBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(changes).To(Equal("Updated manifest with release1-release original-release1-version, release2-release original-release2-version"))
	})

	Context("failure cases", func() {
		It("ensures there is a releases key at the bottom of the manifest", func() {
			releases := []string{"release1", "release2"}
			badManifest := []byte(`
name:
stemcells:
other_key:
`)
			_, _, err := manifest.UpdateReleases(releases, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError("releases was not found at the bottom of the manifest"))
		})

		It("ensures there is a stemcell key at the bottom of the manifest", func() {
			releases := []string{"release1", "release2"}
			badManifest := []byte(`
name:
stemcells:
releases:
other_key:
`)
			_, _, err := manifest.UpdateReleases(releases, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError("stemcells was not found at the bottom of the manifest"))
		})

		It("returns an error when there are keys other than release and stemcells", func() {
			releases := []string{"release1", "release2"}
			badManifest := []byte(`
name:
releases:
stemcells:
other_key:
`)
			_, _, err := manifest.UpdateReleases(releases, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(`found keys other than "releases" and "stemcells" at the bottom of the manifest`))
		})

		It("returns errors instead of panicking when url is missing", func() {
			releases := []string{"missing-url"}

			_, _, err := manifest.UpdateReleases(releases, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-url-release/url: no such file or directory"))
		})

		It("returns errors instead of panicking when version is missing", func() {
			releases := []string{"missing-version"}

			_, _, err := manifest.UpdateReleases(releases, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-version-release/version: no such file or directory"))
		})

		It("returns errors instead of panicking when sha1 is missing", func() {
			releases := []string{"missing-sha1"}

			_, _, err := manifest.UpdateReleases(releases, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-sha1-release/sha1: no such file or directory"))
		})

		It("returns an error when the manifest is not valid yaml", func() {
			releases := []string{"release1", "release2"}
			cfDeploymentManifest := []byte(`
%%%
releases:
%%%
`)
			_, _, err := manifest.UpdateReleases(releases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(ContainSubstring("could not find expected directive name")))
		})

		It("returns an error when the releases section is malformed", func() {
			releases := []string{"release1", "release2"}
			cfDeploymentManifest := []byte(`
name: my-deployment
releases:
- wrong type
stemcells:
- alias: my-stemcell
`)

			_, _, err := manifest.UpdateReleases(releases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(ContainSubstring("`wrong type` into common.Release")))
		})

		It("returns an error when the yaml marshaller fails", func() {
			marshalFailFunc := func(interface{}) ([]byte, error) {
				return nil, errors.New("failed to marshal yaml")
			}
			releases := []string{"release1", "release2"}

			_, _, err := manifest.UpdateReleases(releases, goodBuildDir, cfDeploymentManifest, marshalFailFunc, yaml.Unmarshal)
			Expect(err).To(MatchError("failed to marshal yaml"))
		})

		It("returns an error when the yaml unmarshaller fails", func() {
			unmarshalFailFunc := func([]byte, interface{}) error {
				return errors.New("failed to unmarshal yaml")
			}
			releases := []string{"release1", "release2"}

			_, _, err := manifest.UpdateReleases(releases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, unmarshalFailFunc)
			Expect(err).To(MatchError("failed to unmarshal yaml"))
		})
	})
})

var _ = Describe("UpdateStemcell", func() {
	var (
		brokenBuildDir string
		goodBuildDir   string

		cfDeploymentManifest []byte
	)

	BeforeEach(func() {
		brokenBuildDir = "../fixtures/broken-build"
		goodBuildDir = "../fixtures/build"

		var err error

		cfDeploymentManifest, err = os.ReadFile("../fixtures/cf-deployment.yml")
		Expect(err).NotTo(HaveOccurred())
	})

	It("updates the stemcell without modifying the rest and returns the list of changes when the version changes", func() {
		updatedStemcellFixture, err := os.ReadFile("../fixtures/updated_version_stemcell.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateStemcell([]string{}, "../fixtures/build-with-updated-stemcell-version", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^stemcells:$)`)
		updatedManifestStemcellIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestStemcellIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestStemcellIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestStemcellIndex]

		updatedManifestStemcell := updatedManifest[updatedManifestStemcellIndex:]

		Expect(string(cfDeploymentPreamble)).To(Equal(string(updatedManifestPreamble)), "the preamble was changed by running the program")
		Expect(string(updatedManifestStemcell)).To(Equal(string(updatedStemcellFixture)))

		Expect(changes).To(Equal("Updated manifest with ubuntu-trusty stemcell updated-stemcell-version"))
	})

	It("takes stemcell OS from the stemcell input when the OS is different from the one in the base manifest", func() {
		updatedStemcellFixture, err := os.ReadFile("../fixtures/updated_stemcell_os_and_releases.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateStemcell(nil, "../fixtures/build-with-different-stemcell-os", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

		Expect(err).NotTo(HaveOccurred())
		Expect(updatedManifest).To(MatchYAML(updatedStemcellFixture))
		Expect(changes).To(Equal("Updated manifest with ubuntu-foo stemcell 0.1"))
	})

	Context("failure cases", func() {
		Context("when there is not a stemcell key at the bottom of the manifest", func() {
			It("returns an error", func() {
				badManifest := []byte(`
name:
stemcells:
releases:
other_key:
`)
				_, _, err := manifest.UpdateStemcell([]string{}, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
				Expect(err).To(MatchError("stemcells was not found at the bottom of the manifest"))
			})
		})

		Context("when there are keys other than release and stemcells", func() {
			It("returns an error", func() {
				badManifest := []byte(`
name:
releases:
stemcells:
other_key:
`)
				_, _, err := manifest.UpdateStemcell([]string{}, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
				Expect(err).To(MatchError(`found keys other than "releases" and "stemcells" at the bottom of the manifest`))
			})
		})

		Context("when the stemcell version is missing", func() {
			BeforeEach(func() {
				brokenBuildDir = "../fixtures/broken-build/missing-version-stemcell"
			})

			It("returns an error", func() {
				_, _, err := manifest.UpdateStemcell([]string{}, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

				Expect(err.Error()).To(Equal("open ../fixtures/broken-build/missing-version-stemcell/stemcell/version: no such file or directory"))
			})
		})

		Context("when the stemcell url is missing", func() {
			BeforeEach(func() {
				brokenBuildDir = "../fixtures/broken-build/missing-url-stemcell"
			})

			It("returns an error", func() {
				_, _, err := manifest.UpdateStemcell([]string{}, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

				Expect(err.Error()).To(Equal("open ../fixtures/broken-build/missing-url-stemcell/stemcell/url: no such file or directory"))
			})
		})

		Context("when the stemcell OS cannot be determined from the URL", func() {
			BeforeEach(func() {
				brokenBuildDir = "../fixtures/broken-build/bad-url-stemcell"
			})

			It("returns an error", func() {
				_, _, err := manifest.UpdateStemcell([]string{}, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

				Expect(err.Error()).To(Equal("Stemcell URL does not contain 'ubuntu': bad-stemcell-url"))
			})
		})

		Context("when the manifest is not valid yaml", func() {
			It("returns an error", func() {
				cfDeploymentManifest := []byte(`
%%%
releases:
%%%
`)
				_, _, err := manifest.UpdateStemcell([]string{}, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
				Expect(err).To(MatchError(ContainSubstring("could not find expected directive name")))
			})
		})

		Context("when the stemcells section is malformed", func() {
			It("returns an error", func() {
				cfDeploymentManifest := []byte(`
name: my-deployment
releases:
- name: my-release
stemcells:
- wrong type
`)

				_, _, err := manifest.UpdateStemcell([]string{}, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
				Expect(err).To(MatchError(ContainSubstring("`wrong type` into manifest.Stemcell")))
			})
		})
	})
})
