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
				{
					Name: "product-with-hyphens",
					SHA1: "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
					Stemcell: manifest.Stemcell{
						OS:      "some-stemcell",
						Version: "1.2",
					},
					Version: "1.2.3",
					URL:     fmt.Sprintf("https://storage.googleapis.com/cf-deployment-compiled-releases/%s", hyphenPath),
				},
				{
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

	Describe("Write", func() {
		var (
			actualError error
		)

		JustBeforeEach(func() {
			actualError = opsfile.Write()
		})

		Context("when the opsfile has a filled array of Ops", func() {
			BeforeEach(func() {
				opsfile.Ops = []Op{
					{
						Type: "replace",
						Path: "/releases/name=some-buildpack",
						Value: manifest.Release{
							Name: "some-buildpack",
							SHA1: "123456",
							Stemcell: manifest.Stemcell{
								OS:      "some-stemcell",
								Version: "1.2",
							},
							Version: "1.2.3",
							URL:     "some-url/some-buildpack.com",
						},
					},
					{
						Type: "replace",
						Path: "/releases/name=some-component",
						Value: manifest.Release{
							Name: "some-component",
							SHA1: "aabbff",
							Stemcell: manifest.Stemcell{
								OS:      "some-stemcell",
								Version: "1.2",
							},
							Version: "4.5.6",
							URL:     "some-url/some-component.com",
						},
					},
				}
			})

			It("generates the opsfile for compiled releases with the updated stemcell", func() {
				Expect(actualError).NotTo(HaveOccurred())

				actualContents, err := ioutil.ReadFile(opsfileOutPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(actualContents)).To(Equal(`## GENERATED FILE. DO NOT EDIT
---
- type: replace
  path: /releases/name=some-buildpack
  value:
    name: some-buildpack
    sha1: "123456"
    stemcell:
      os: some-stemcell
      version: "1.2"
    version: 1.2.3
    url: some-url/some-buildpack.com
- type: replace
  path: /releases/name=some-component
  value:
    name: some-component
    sha1: aabbff
    stemcell:
      os: some-stemcell
      version: "1.2"
    version: 4.5.6
    url: some-url/some-component.com
`))
			})
		})
	})
})
