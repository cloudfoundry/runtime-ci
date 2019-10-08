package compiledrelease

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OpsfileUpdater", func() {
	var (
		buildDir string

		opsfileOutPath string
		opsfileUpdater *OpsfileUpdater
	)

	BeforeEach(func() {
		var err error
		buildDir, err = ioutil.TempDir("", "concourseio-rootdir-")
		Expect(err).ToNot(HaveOccurred())

		opsfileOutPath = filepath.Join(buildDir, "ops-file.yml")

		opsfileUpdater = NewOpsfileUpdater(buildDir, opsfileOutPath)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(buildDir)).To(Succeed())
	})

	Describe("Load", func() {
		var (
			hyphenPath     string
			singlewordPath string

			actualError error
		)

		JustBeforeEach(func() {
			actualError = opsfileUpdater.Load()
		})

		Context("when there are releases", func() {
			BeforeEach(func() {
				hyphenPath = "product-with-hyphens-1.2.3-some-stemcell-1.2-00000000-000000-000000000.tgz"
				singlewordPath = "singleword-4.5.6-some-stemcell-1.2-00000000-000000-000000000.tgz"

				Expect(ioutil.WriteFile(filepath.Join(buildDir, hyphenPath), []byte("hello world"), 0777)).To(Succeed())
				Expect(ioutil.WriteFile(filepath.Join(buildDir, singlewordPath), []byte("hello kitty"), 0777)).To(Succeed())
			})

			It("load a list of Releases from the compiled-releases directory", func() {
				Expect(actualError).NotTo(HaveOccurred())

				expectedReleases := []bosh.Release{
					{
						Name: "product-with-hyphens",
						SHA1: "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
						Stemcell: bosh.Stemcell{
							OS:      "some-stemcell",
							Version: "1.2",
						},
						Version: "1.2.3",
						URL:     fmt.Sprintf("https://storage.googleapis.com/cf-deployment-compiled-releases/%s", hyphenPath),
					},
					{
						Name: "singleword",
						SHA1: "89f53c408c8bd119b92a295f30963de7dcb00f2f",
						Stemcell: bosh.Stemcell{
							OS:      "some-stemcell",
							Version: "1.2",
						},
						Version: "4.5.6",
						URL:     fmt.Sprintf("https://storage.googleapis.com/cf-deployment-compiled-releases/%s", singlewordPath),
					},
				}

				Expect(opsfileUpdater.releases).To(ConsistOf(expectedReleases))
			})
		})

		Context("when there are no releases", func() {
			It("will return an error", func() {
				Expect(actualError).To(MatchError(&NoReleasesErr{}))
			})
		})

		Context("when there is an invalid release", func() {
			BeforeEach(func() {
				singlewordPath = "invalid-1.2.3.tgz"

				Expect(ioutil.WriteFile(filepath.Join(buildDir, singlewordPath), []byte("hello kitty"), 0777)).To(Succeed())
			})

			It("will return an error", func() {
				Expect(actualError).To(MatchError("invalid tarball name syntax: invalid-1.2.3.tgz"))
			})
		})
	})

	Describe("Update", func() {
		var (
			stemcellArg bosh.Stemcell

			actualError error
		)

		JustBeforeEach(func() {
			actualError = opsfileUpdater.Update(stemcellArg)
		})

		Context("when the opsfileUpdater has releases", func() {
			BeforeEach(func() {
				opsfileUpdater.releases = []bosh.Release{
					{
						Name: "some-buildpack",
						SHA1: "123456",
						Stemcell: bosh.Stemcell{
							OS:      "some-stemcell",
							Version: "1.2",
						},
						Version: "1.2.3",
						URL:     "some-url/some-buildpack.com",
					},
					{
						Name: "some-component",
						SHA1: "aabbff",
						Stemcell: bosh.Stemcell{
							OS:      "some-stemcell",
							Version: "1.2",
						},
						Version: "4.5.6",
						URL:     "some-url/some-component.com",
					},
				}

				stemcellArg = bosh.Stemcell{
					OS:      "some-stemcell",
					Version: "1.2",
				}
			})

			It("populates the ops in the opsfileUpdater", func() {
				Expect(actualError).ToNot(HaveOccurred())

				Expect(opsfileUpdater.ops).To(ConsistOf(
					Op{
						Type: "replace",
						Path: "/releases/name=some-buildpack",
						Value: bosh.Release{
							Name: "some-buildpack",
							SHA1: "123456",
							Stemcell: bosh.Stemcell{
								OS:      "some-stemcell",
								Version: "1.2",
							},
							Version: "1.2.3",
							URL:     "some-url/some-buildpack.com",
						},
					},
					Op{
						Type: "replace",
						Path: "/releases/name=some-component",
						Value: bosh.Release{
							Name: "some-component",
							SHA1: "aabbff",
							Stemcell: bosh.Stemcell{
								OS:      "some-stemcell",
								Version: "1.2",
							},
							Version: "4.5.6",
							URL:     "some-url/some-component.com",
						},
					},
				))
			})

			Context("when the stemcell does not match", func() {
				BeforeEach(func() {
					stemcellArg = bosh.Stemcell{
						OS:      "some-stemcell",
						Version: "3.4",
					}
				})

				It("will return an error", func() {
					Expect(actualError).To(MatchError("stemcell mismatch"))
				})
			})
		})

		Context("when there are no releases", func() {
			It("will return an error", func() {
				Expect(actualError).To(MatchError(&NoReleasesErr{}))
			})
		})
	})

	Describe("Write", func() {
		var (
			actualError error
		)

		JustBeforeEach(func() {
			actualError = opsfileUpdater.Write()
		})

		Context("when there are somehow no ops from releases", func() {
			It("will return an error and not write to the Outfile", func() {
				Expect(actualError).To(MatchError(&NoReleasesErr{}))

				_, err := ioutil.ReadFile(opsfileOutPath)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})

		Context("when the opsfileUpdater has a filled array of ops", func() {
			BeforeEach(func() {
				opsfileUpdater.ops = []Op{
					{
						Type: "replace",
						Path: "/releases/name=some-buildpack",
						Value: bosh.Release{
							Name: "some-buildpack",
							SHA1: "123456",
							Stemcell: bosh.Stemcell{
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
						Value: bosh.Release{
							Name: "some-component",
							SHA1: "aabbff",
							Stemcell: bosh.Stemcell{
								OS:      "some-stemcell",
								Version: "1.2",
							},
							Version: "4.5.6",
							URL:     "some-url/some-component.com",
						},
					},
				}
			})

			It("generates the opsfileUpdater for compiled releases with the updated stemcell", func() {
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
    url: some-url/some-buildpack.com
    version: 1.2.3
- type: replace
  path: /releases/name=some-component
  value:
    name: some-component
    sha1: aabbff
    stemcell:
      os: some-stemcell
      version: "1.2"
    url: some-url/some-component.com
    version: 4.5.6
`))
			})
		})
	})
})
