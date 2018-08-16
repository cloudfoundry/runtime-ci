package common_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/runtime-ci/util/update-manifest-releases/common"
)
var _ = Describe("Common", func() {
	var (
		buildDir          string
	)

	BeforeEach(func() {
		buildDir = "../fixtures/broken-build"
	})

	Context("#GetReleaseFromFile", func() {
		Context("when release folder has all required files", func() {
			It("returns release desired release from the build dir", func() {
				release, err := common.GetReleaseFromFile(buildDir, "good-release")

				Expect(err).NotTo(HaveOccurred())
				Expect(release.Name).To(Equal("good-release"))
				Expect(release.URL).To(Equal("https://download.com/release1"))
				Expect(release.SHA1).To(Equal("XXXXXXXXXXXXXX"))
				Expect(release.Version).To(Equal("1.1"))
			})
		})

		Context("when release folder is missing files", func() {
			It("doesn't error when both sha1 and url are missing", func() {
				release, err := common.GetReleaseFromFile(buildDir, "missing-url-and-sha1")

				Expect(err).NotTo(HaveOccurred())
				Expect(release.Name).To(Equal("missing-url-and-sha1"))
				Expect(release.Version).To(Equal("updated-missing-url-and-sha1-version"))
			})

			It("errors when sha1 is missing", func() {
				_, err := common.GetReleaseFromFile(buildDir, "missing-sha1")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("open ../fixtures/broken-build/missing-sha1-release/sha1: no such file or directory"))
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

	Context("#StemcellInfoFromTarballName", func() {
		It("parses a stemcell version from a valid release tarball names", func() {
			// release and stemcell name with no dashes
			stemcellVersion, stemcellName, err := common.StemcellInfoFromTarballName("release-1.0-stemcell-name-2.0-4-5-6.tgz", "release", "1.0")
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellVersion).To(Equal("2.0"))
			Expect(stemcellName).To(Equal("stemcell-name"))

			// release and stemcell name with one dashe
			stemcellVersion, stemcellName, err = common.StemcellInfoFromTarballName("release-name-1.0-stemcell-name-2.0-45-23-44.tgz", "release-name",  "1.0")
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellVersion).To(Equal("2.0"))
			Expect(stemcellName).To(Equal("stemcell-name"))

			// release and stemcell name with many dashes
			stemcellVersion, stemcellName, err = common.StemcellInfoFromTarballName("release-name-long-1.0-stemcell-name-long-2.0-4-5-6.tgz", "release-name-long",  "1.0")
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellVersion).To(Equal("2.0"))
			Expect(stemcellName).To(Equal("stemcell-name-long"))

			// stemcell version with no dots
			stemcellVersion, _, err = common.StemcellInfoFromTarballName("release-name-1.0-stemcell-name-2-4-5-6.tgz", "release-name", "1.0")
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellVersion).To(Equal("2"))

			// stemcell version with many dots
			stemcellVersion, _, err = common.StemcellInfoFromTarballName("release-name-1.0-stemcell-name-2.1.3-4-5-6.tgz", "release-name", "1.0")
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellVersion).To(Equal("2.1.3"))
		})

		It("should fail if the tarball name does not start with the release name or release version", func() {
			_, _, err := common.StemcellInfoFromTarballName("release-1.0-stemcell-2.0-4-5-6.tgz", "different-release", "1.0")
			Expect(err).To(HaveOccurred())

			_, _, err = common.StemcellInfoFromTarballName("release-1.0-stemcell-2.0-4-5-6.tgz", "release", "1000000.0")
			Expect(err).To(HaveOccurred())
		})

		It("should fail if the tarball name has an incorrect number of timestamp words", func() {
			_, _, err := common.StemcellInfoFromTarballName("release-1.0-stemcell-2-3-4.tgz", "release", "1.0")
			Expect(err).To(HaveOccurred())
		})

	})
})