package manifest_test

import (
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
	})

	// TODO: Make CF Deployment manifest parameterized so that we can pass it in the real manifest somehow
	It("updates the releases and stemcells to their latest version without modifying the rest", func() {
		releases := []string{"release1", "release2"}
		buildDir := "fixtures/build"

		updatedManifest := manifest.UpdateReleasesAndStemcells(releases, buildDir, cfDeploymentManifest)

		r := regexp.MustCompile(`(?m:^releases:$)`)
		updatedManifestReleasesIndex := r.FindSubmatchIndex([]byte(updatedManifest))[0]
		cfDeploymentManifestReleasesIndex := r.FindSubmatchIndex([]byte(cfDeploymentManifest))[0]
		cfDeploymentPreamble := cfDeploymentManifest[:cfDeploymentManifestReleasesIndex]
		updatedManifestPreamble := updatedManifest[:updatedManifestReleasesIndex]

		updatedManifestReleasesAndStemcells := updatedManifest[updatedManifestReleasesIndex:]

		Expect(cfDeploymentPreamble).To(MatchYAML(updatedManifestPreamble))
		Expect(updatedManifestReleasesAndStemcells).To(MatchYAML(updatedReleasesAndStemcellsFixture))
	})
})
