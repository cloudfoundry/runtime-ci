package common_test

import (
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common", func() {
	var (
		buildDir string
	)

	BeforeEach(func() {
		buildDir = "../fixtures/broken-build"
	})

	Context("#GetReleaseFromFile", func() {
		Context("when release folder has all required files", func() {
			It("returns the desired release from the build dir", func() {
				release, err := common.GetReleaseFromFile(buildDir, "good-release")

				Expect(err).NotTo(HaveOccurred())
				Expect(release.Name).To(Equal("good-release"))
				Expect(release.URL).To(Equal("https://download.com/release1"))
				Expect(release.SHA1).To(Equal("sha256:XXXXXXXXXXXXXX"))
				Expect(release.Version).To(Equal("1.1"))
			})
		})

		Context("when the release folder is from the github-release-resource", func() {
			It("returns the desired release from the build dir", func() {
				release, err := common.GetReleaseFromFile(buildDir, "good-github-release")

				Expect(err).NotTo(HaveOccurred())
				Expect(release.Name).To(Equal("good-github-release"))
				Expect(release.URL).To(BeEmpty())
				Expect(release.SHA1).To(BeEmpty())
				Expect(release.Version).To(Equal("1.2.3"))
			})
		})

		Context("when release folder is missing files", func() {
			It("errors when sha256 is missing", func() {
				_, err := common.GetReleaseFromFile(buildDir, "missing-sha256")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("open ../fixtures/broken-build/missing-sha256-release/sha256: no such file or directory"))
			})

			It("errors when url is missing", func() {
				_, err := common.GetReleaseFromFile(buildDir, "missing-url")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("open ../fixtures/broken-build/missing-url-release/url: no such file or directory"))
			})

			It("errors when version is missing", func() {
				_, err := common.GetReleaseFromFile(buildDir, "missing-version")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("open ../fixtures/broken-build/missing-version-release/version: no such file or directory"))
			})
		})
	})

	Context("#InfoFromTarballName", func() {
		It("parses info from a valid release tarball name with no dashes", func() {
			releaseVersion, stemcellVersion, stemcellName, err := common.InfoFromTarballName("release-1.0-stemcell-2.0-4-5-6.tgz", "release")
			Expect(err).NotTo(HaveOccurred())
			Expect(releaseVersion).To(Equal("1.0"))
			Expect(stemcellVersion).To(Equal("2.0"))
			Expect(stemcellName).To(Equal("stemcell"))
		})

		It("parses info from a valid release tarball name with one dash", func() {
			releaseVersion, stemcellVersion, stemcellName, err := common.InfoFromTarballName("release-name-1.0-stemcell-name-2.0-45-23-44.tgz", "release-name")
			Expect(err).NotTo(HaveOccurred())
			Expect(releaseVersion).To(Equal("1.0"))
			Expect(stemcellVersion).To(Equal("2.0"))
			Expect(stemcellName).To(Equal("stemcell-name"))
		})

		It("parses info from a valid release tarball name with multiple dashes", func() {
			releaseVersion, stemcellVersion, stemcellName, err := common.InfoFromTarballName("release-name-long-1.0-stemcell-name-long-2.0-4-5-6.tgz", "release-name-long")
			Expect(err).NotTo(HaveOccurred())
			Expect(releaseVersion).To(Equal("1.0"))
			Expect(stemcellVersion).To(Equal("2.0"))
			Expect(stemcellName).To(Equal("stemcell-name-long"))
		})

		It("parses info from a valid release tarball name with a stemcell version with no dots", func() {
			releaseVersion, stemcellVersion, stemcellName, err := common.InfoFromTarballName("release-name-1.0-stemcell-name-2-4-5-6.tgz", "release-name")
			Expect(err).NotTo(HaveOccurred())
			Expect(releaseVersion).To(Equal("1.0"))
			Expect(stemcellVersion).To(Equal("2"))
			Expect(stemcellName).To(Equal("stemcell-name"))
		})

		It("parses info from a valid release tarball name with a stemcell version with multiple dots", func() {
			releaseVersion, stemcellVersion, stemcellName, err := common.InfoFromTarballName("release-name-1.0-stemcell-name-2.1.3-4-5-6.tgz", "release-name")
			Expect(err).NotTo(HaveOccurred())
			Expect(releaseVersion).To(Equal("1.0"))
			Expect(stemcellVersion).To(Equal("2.1.3"))
			Expect(stemcellName).To(Equal("stemcell-name"))
		})

		It("should fail if the tarball name does not start with the release name", func() {
			_, _, _, err := common.InfoFromTarballName("release-1.0-stemcell-2.0-4-5-6.tgz", "different-release")
			Expect(err.Error()).To(ContainSubstring("invalid tarball name syntax"))
		})

		It("should fail if the tarball name has an incorrect number of timestamp words", func() {
			_, _, _, err := common.InfoFromTarballName("release-1.0-stemcell-2-3-4.tgz", "release")
			Expect(err.Error()).To(ContainSubstring("invalid tarball name syntax"))
		})

	})
})
