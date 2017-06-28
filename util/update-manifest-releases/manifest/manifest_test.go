package manifest_test

import (
	"errors"
	"io/ioutil"
	"regexp"

	yaml "gopkg.in/yaml.v2"

	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/manifest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateReleasesAndStemcells", func() {
	var (
		brokenBuildDir    string
		goodBuildDir      string
		noChangesBuildDir string

		updatedReleasesAndStemcellsFixture []byte
		cfDeploymentManifest               []byte
	)

	BeforeEach(func() {
		brokenBuildDir = "../fixtures/broken-build"
		goodBuildDir = "../fixtures/build"
		noChangesBuildDir = "../fixtures/nochanges-build"

		var err error

		updatedReleasesAndStemcellsFixture, err = ioutil.ReadFile("../fixtures/updated_releases_and_stemcells.yml")
		Expect(err).NotTo(HaveOccurred())

		cfDeploymentManifest, err = ioutil.ReadFile("../fixtures/cf-deployment.yml")
		Expect(err).NotTo(HaveOccurred())
	})

	It("when the sha changes, updates the releases and stemcells without modifying the rest and returns the list of changes", func() {
		releases := []string{"release1", "release2"}
		updatedReleasesAndStemcellsFixture, err := ioutil.ReadFile("../fixtures/updated_sha_releases_and_stemcells.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateReleasesAndStemcells(releases, "../fixtures/build-with-updated-sha", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleasesAndStemcells := updatedManifest[updatedManifestReleasesIndex:]

		Expect(string(cfDeploymentPreamble)).To(Equal(string(updatedManifestPreamble)), "the preamble was changed by running the program")
		Expect(string(updatedManifestReleasesAndStemcells)).To(Equal(string(updatedReleasesAndStemcellsFixture)))

		Expect(changes).To(Equal("Updated manifest with release2-release original-release2-version, ubuntu-trusty stemcell updated-stemcell-version"))
	})

	It("when the version changes, updates the releases and stemcells without modifying the rest and returns the list of changes", func() {
		releases := []string{"release1", "release2"}
		updatedReleasesAndStemcellsFixture, err := ioutil.ReadFile("../fixtures/updated_version_releases_and_stemcells.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateReleasesAndStemcells(releases, "../fixtures/build-with-updated-version", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleasesAndStemcells := updatedManifest[updatedManifestReleasesIndex:]

		Expect(string(cfDeploymentPreamble)).To(Equal(string(updatedManifestPreamble)), "the preamble was changed by running the program")
		Expect(string(updatedManifestReleasesAndStemcells)).To(Equal(string(updatedReleasesAndStemcellsFixture)))

		Expect(changes).To(Equal("Updated manifest with release2-release updated-release2-version, ubuntu-trusty stemcell updated-stemcell-version"))
	})

	It("when the url changes, updates the releases and stemcells without modifying the rest and returns the list of changes", func() {
		releases := []string{"release1", "release2"}
		updatedReleasesAndStemcellsFixture, err := ioutil.ReadFile("../fixtures/updated_url_releases_and_stemcells.yml")
		Expect(err).NotTo(HaveOccurred())

		updatedManifest, changes, err := manifest.UpdateReleasesAndStemcells(releases, "../fixtures/build-with-updated-url", cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleasesAndStemcells := updatedManifest[updatedManifestReleasesIndex:]

		Expect(string(cfDeploymentPreamble)).To(Equal(string(updatedManifestPreamble)), "the preamble was changed by running the program")
		Expect(string(updatedManifestReleasesAndStemcells)).To(Equal(string(updatedReleasesAndStemcellsFixture)))

		Expect(changes).To(Equal("Updated manifest with release2-release original-release2-version, ubuntu-trusty stemcell updated-stemcell-version"))
	})

	It("provides a default commit message if no version updates were performed", func() {
		releases := []string{"release1", "release2"}
		_, changes, err := manifest.UpdateReleasesAndStemcells(releases, noChangesBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(changes).To(Equal("No manifest release or stemcell version updates"))
	})

	It("when there exist update releases that are not in the manifest releases, it adds them to resulting list of releases", func() {
		updateReleases := []string{"release1"}
		cfDeploymentManifest := []byte(`
name: my-deployment
releases:
  - name: fooRelease
stemcells:
`)
		resultingManifest, _, err := manifest.UpdateReleasesAndStemcells(updateReleases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).ToNot(HaveOccurred())

		var releasesAndStemcells manifest.Manifest
		err = yaml.Unmarshal(resultingManifest, &releasesAndStemcells)
		Expect(err).ToNot(HaveOccurred())

		Expect(releasesAndStemcells.Releases).To(ContainElement(common.Release{
			Name:    "release1",
			URL:     "original-release1-url",
			Version: "original-release1-version",
			SHA1:    "original-release1-sha1",
		}))

	})

	It("adds all the releases and stemcells to the commit message if no releases exist", func() {
		releases := []string{"release1", "release2"}
		cfDeploymentManifest := []byte(`
name: my-deployment
releases:
stemcells:
`)

		_, changes, err := manifest.UpdateReleasesAndStemcells(releases, noChangesBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
		Expect(err).NotTo(HaveOccurred())

		Expect(changes).To(Equal("Updated manifest with release1-release original-release1-version, release2-release original-release2-version, ubuntu-trusty stemcell original-stemcell-version"))
	})

	Context("failure cases", func() {
		It("ensures there is a releases key at the bottom of the manifest", func() {
			releases := []string{"release1", "release2"}
			badManifest := []byte(`
name:
stemcells:
other_key:
`)
			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
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
			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
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
			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, badManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(`found keys other than "releases" and "stemcells" at the bottom of the manifest`))
		})

		It("returns errors instead of panicking when url is missing", func() {
			releases := []string{"missing-url"}

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-url-release/url: no such file or directory"))
		})

		It("returns errors instead of panicking when version is missing", func() {
			releases := []string{"missing-version"}

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-version-release/version: no such file or directory"))
		})

		It("returns errors instead of panicking when sha1 is missing", func() {
			releases := []string{"missing-sha1"}

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/missing-sha1-release/sha1: no such file or directory"))
		})

		It("returns errors instead of panicking when sha1 is missing", func() {
			releases := []string{"good-release"}

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, brokenBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)

			Expect(err).To(MatchError("open ../fixtures/broken-build/stemcell/version: no such file or directory"))
		})

		It("returns an error when the manifest is not valid yaml", func() {
			releases := []string{"release1", "release2"}
			cfDeploymentManifest := []byte(`
%%%
releases:
%%%
`)
			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
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

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(ContainSubstring("`wrong type` into common.Release")))
		})

		It("returns an error when the stemcells section is malformed", func() {
			releases := []string{"release1", "release2"}
			cfDeploymentManifest := []byte(`
name: my-deployment
releases:
- name: my-release
stemcells:
- wrong type
`)

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, yaml.Unmarshal)
			Expect(err).To(MatchError(ContainSubstring("`wrong type` into manifest.Stemcell")))
		})

		It("returns an error when the yaml marshaller fails", func() {
			marshalFailFunc := func(interface{}) ([]byte, error) {
				return nil, errors.New("failed to marshal yaml")
			}
			releases := []string{"release1", "release2"}

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, cfDeploymentManifest, marshalFailFunc, yaml.Unmarshal)
			Expect(err).To(MatchError("failed to marshal yaml"))
		})

		It("returns an error when the yaml unmarshaller fails", func() {
			unmarshalFailFunc := func([]byte, interface{}) error {
				return errors.New("failed to unmarshal yaml")
			}
			releases := []string{"release1", "release2"}

			_, _, err := manifest.UpdateReleasesAndStemcells(releases, goodBuildDir, cfDeploymentManifest, yaml.Marshal, unmarshalFailFunc)
			Expect(err).To(MatchError("failed to unmarshal yaml"))
		})
	})
})
