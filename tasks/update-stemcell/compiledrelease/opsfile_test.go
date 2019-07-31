package compiledrelease

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Opsfile", func() {
	var (
		buildDir string

		opsfileOutPath string
		opsfile        *Opsfile
	)

	BeforeEach(func() {
		var err error
		buildDir, err = ioutil.TempDir("", "concourseio-rootdir-")
		Expect(err).ToNot(HaveOccurred())

		opsfileOutPath = filepath.Join(buildDir, "ops-file.yml")

		opsfile = NewOpsfile(buildDir, opsfileOutPath)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(buildDir)).To(Succeed())
	})

	Describe("Load", func() {
		var (
			hyphenPath     string
			singlewordPath string
		)

		BeforeEach(func() {
			hyphenPath = "product-with-hyphens-1.2.3-some-stemcell-1.2-00000000-000000-000000000.tgz"
			singlewordPath = "singleword-4.5.6-some-stemcell-1.2-00000000-000000-000000000.tgz"

			Expect(ioutil.WriteFile(filepath.Join(buildDir, hyphenPath), []byte("hello world"), 0777)).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(buildDir, singlewordPath), []byte("hello kitty"), 0777)).To(Succeed())
		})

		JustBeforeEach(func() {
			err := opsfile.Load()
			Expect(err).NotTo(HaveOccurred())
		})

		It("load a list of Releases from the compiled-releases directory", func() {
			expectedReleases := []manifest.Release{
				manifest.Release{
					Name: "product-with-hyphens",
					SHA1: "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
					Stemcell: manifest.Stemcell{
						OS:      "some-stemcell",
						Version: "1.2",
					},
					Version: "1.2.3",
					URL:     fmt.Sprintf("https://storage.googleapis.com/cf-deployment-compiled-releases/%s", hyphenPath),
				},
				manifest.Release{
					Name: "singleword",
					SHA1: "89f53c408c8bd119b92a295f30963de7dcb00f2f",
					Stemcell: manifest.Stemcell{
						OS:      "some-stemcell",
						Version: "1.2",
					},
					Version: "4.5.6",
					URL:     fmt.Sprintf("https://storage.googleapis.com/cf-deployment-compiled-releases/%s", singlewordPath),
				},
			}
			Expect(opsfile.releases).To(ConsistOf(expectedReleases))
		})
	})
})
