package manifest_test

import (
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/create-binaries-manifest-section/manifest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateReleasesAndStemcells", func() {
	var (
		updatedReleasesAndStemcellsFixture []byte
		cfDeploymentManifest               []byte
	)

	BeforeEach(func() {
		var err error

		updatedReleasesAndStemcellsFixture, err = ioutil.ReadFile("fixtures/updated_releases_and_stemcells.yml")
		Expect(err).NotTo(HaveOccurred())

		cfDeploymentManifest, err = ioutil.ReadFile("fixtures/cf-deployment.yml")
		Expect(err).NotTo(HaveOccurred())

		manifest.ResetYAMLMarshal()
	})

	// TODO: Make CF Deployment manifest parameterized so that we can pass it in the real manifest somehow
	It("updates the releases and stemcells to their latest version without modifying the rest", func() {
		releases := []string{"release1", "release2"}
		buildDir := "fixtures/build"

		updatedManifest, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)
		Expect(err).NotTo(HaveOccurred())

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleasesAndStemcells := updatedManifest[updatedManifestReleasesIndex:]

		Expect(cfDeploymentPreamble).To(MatchYAML(updatedManifestPreamble))
		Expect(updatedManifestReleasesAndStemcells).To(MatchYAML(updatedReleasesAndStemcellsFixture))
	})

	Context("failure cases", func() {
		It("returns errors instead of panicking when url is missing", func() {
			releases := []string{"missing-url"}
			buildDir := "fixtures/broken-build"

			_, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)

			Expect(err).To(MatchError("open fixtures/broken-build/missing-url-release/url: no such file or directory"))
		})

		It("returns errors instead of panicking when version is missing", func() {
			releases := []string{"missing-version"}
			buildDir := "fixtures/broken-build"

			_, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)

			Expect(err).To(MatchError("open fixtures/broken-build/missing-version-release/version: no such file or directory"))
		})

		It("returns errors instead of panicking when sha1 is missing", func() {
			releases := []string{"missing-sha1"}
			buildDir := "fixtures/broken-build"

			_, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)

			Expect(err).To(MatchError("open fixtures/broken-build/missing-sha1-release/sha1: no such file or directory"))
		})

		It("returns errors instead of panicking when sha1 is missing", func() {
			releases := []string{"good-release"}
			buildDir := "fixtures/broken-build"

			_, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)

			Expect(err).To(MatchError("open fixtures/broken-build/stemcell/version: no such file or directory"))
		})

		It("returns an error when the yaml marshaller fails", func() {
			manifest.SetYAMLMarshal(func(interface{}) ([]byte, error) {
				return nil, errors.New("failed to marshal yaml")
			})
			releases := []string{"release1", "release2"}
			buildDir := "fixtures/build"

			_, err := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)
			Expect(err).To(MatchError("failed to marshal yaml"))
		})
	})
})
