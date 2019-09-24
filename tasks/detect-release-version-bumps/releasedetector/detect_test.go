package releasedetector

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Detect", func() {
	var (
		prevManifestPath string
		nextManifestPath string

		releaseDetector *ReleaseDetector

		err error
	)

	BeforeEach(func() {
		buildDir, err := ioutil.TempDir("", "concourseio-rootdir-")
		Expect(err).ToNot(HaveOccurred())
		prevManifestPath = filepath.Join(buildDir, "cf-deployment-prev.yml")
		nextManifestPath = filepath.Join(buildDir, "cf-deployment-next.yml")
		releaseDetector = NewReleaseDetector()
	})

	Describe("Load", func() {
		JustBeforeEach(func() {
			err = releaseDetector.Load(prevManifestPath, nextManifestPath)
		})

		Context("if files exist", func() {
			BeforeEach(func() {
				prevDeploymentManifest := []byte(`---
releases:
- name: release-constant
  url: https://release-constant.io/1
  version: 1.0.0
  sha1: 61dc701f2b851f18282cd249aaf153479fdd620c
- name: release-changing
  url: https://release-changing.io/1
  version: 1.0.0
  sha1: b41556af773ea9aec93dd21a9bbf129200849eed
`)
				Expect(ioutil.WriteFile(prevManifestPath, prevDeploymentManifest, 0777)).To(Succeed())

				nextDeploymentManifest := []byte(`---
releases:
- name: release-constant
  url: https://release-constant.io/1
  version: 1.0.0
  sha1: 61dc701f2b851f18282cd249aaf153479fdd620c
- name: release-changing
  url: https://release-changing.io/2
  version: 2.0.0
  sha1: c41556af773ea9aec93dd15052105125910591ef
- name: release-new
  url: https://release-new.io/3
  version: 3.0.0
  sha1: d5435463432ea9aec93dd15052105125910591ef
`)
				Expect(ioutil.WriteFile(nextManifestPath, nextDeploymentManifest, 0777)).To(Succeed())
			})

			It("loads the prev and next cf deployment releases", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(releaseDetector.prevReleases).To(ConsistOf([]bosh.Release{
					{
						Name:    "release-constant",
						SHA1:    "61dc701f2b851f18282cd249aaf153479fdd620c",
						Version: "1.0.0",
						URL:     "https://release-constant.io/1",
					},
					{
						Name:    "release-changing",
						SHA1:    "b41556af773ea9aec93dd21a9bbf129200849eed",
						Version: "1.0.0",
						URL:     "https://release-changing.io/1",
					},
				}))

				Expect(releaseDetector.nextReleases).To(ConsistOf([]bosh.Release{
					{
						Name:    "release-constant",
						SHA1:    "61dc701f2b851f18282cd249aaf153479fdd620c",
						Version: "1.0.0",
						URL:     "https://release-constant.io/1",
					},
					{
						Name:    "release-changing",
						SHA1:    "c41556af773ea9aec93dd15052105125910591ef",
						Version: "2.0.0",
						URL:     "https://release-changing.io/2",
					},
					{
						Name:    "release-new",
						SHA1:    "d5435463432ea9aec93dd15052105125910591ef",
						Version: "3.0.0",
						URL:     "https://release-new.io/3",
					},
				}))
			})
		})

		Context("When one of the files does not exist", func() {
			BeforeEach(func() {
				prevDeploymentManifest := []byte(`---
releases:
- name: release-constant
  url: https://release-constant.io/1
  version: 1.0.0
  sha1: 61dc701f2b851f18282cd249aaf153479fdd620c
- name: release-changing
  url: https://release-changing.io/1
  version: 1.0.0
  sha1: b41556af773ea9aec93dd21a9bbf129200849eed
`)
				Expect(ioutil.WriteFile(prevManifestPath, prevDeploymentManifest, 0777)).To(Succeed())
			})

			It("it gives back the error", func() {
				Expect(err).To(MatchError(fmt.Sprintf("open %s: no such file or directory", nextManifestPath)))
			})
		})
	})

	Describe("DetectUpdatedReleases", func() {
		BeforeEach(func() {
			releaseDetector.prevReleases = []bosh.Release{
				{
					Name:    "release-constant",
					SHA1:    "61dc701f2b851f18282cd249aaf153479fdd620c",
					Version: "1.0.0",
					URL:     "https://release-constant.io",
				},
				{
					Name:    "release-changing",
					SHA1:    "b41556af773ea9aec93dd21a9bbf129200849eed",
					Version: "1.0.0",
					URL:     "https://release-changing.io/1",
				},
			}

			releaseDetector.nextReleases = []bosh.Release{
				{
					Name:    "release-constant",
					SHA1:    "61dc701f2b851f18282cd249aaf153479fdd620c",
					Version: "1.0.0",
					URL:     "https://release-constant.io",
				},
				{
					Name:    "release-changing",
					SHA1:    "c41556af773ea9aec93dd15052105125910591ef",
					Version: "2.0.0",
					URL:     "https://release-changing.io/2",
				},
				{
					Name:    "release-new",
					SHA1:    "d5435463432ea9aec93dd15052105125910591ef",
					Version: "3.0.0",
					URL:     "https://release-new.io/3",
				},
			}
		})

		JustBeforeEach(func() {
			releaseDetector.DetectUpdatedReleases()
		})

		It("generates a list of updated Releases", func() {
			Expect(releaseDetector.updatedReleases).To(ConsistOf([]bosh.Release{
				{
					Name:    "release-changing",
					SHA1:    "c41556af773ea9aec93dd15052105125910591ef",
					Version: "2.0.0",
					URL:     "https://release-changing.io/2",
				},
				{
					Name:    "release-new",
					SHA1:    "d5435463432ea9aec93dd15052105125910591ef",
					Version: "3.0.0",
					URL:     "https://release-new.io/3",
				},
			}))
		})
	})

	Describe("Write", func() {
		var (
			listOfUpdatedReleases []byte
		)

		BeforeEach(func() {
			releaseDetector.updatedReleases = []bosh.Release{
				{
					Name:    "release-changing",
					SHA1:    "c41556af773ea9aec93dd15052105125910591ef",
					Version: "2.0.0",
					URL:     "https://release-changing.io/2",
				},
				{
					Name:    "release-new",
					SHA1:    "d5435463432ea9aec93dd15052105125910591ef",
					Version: "3.0.0",
					URL:     "https://release-new.io/3",
				},
			}
		})

		JustBeforeEach(func() {
			listOfUpdatedReleases, err = releaseDetector.Write()
		})

		It("creates yaml contaning the list of changed releases", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(string(listOfUpdatedReleases)).To(Equal(`releases:
- name: release-changing
  sha1: c41556af773ea9aec93dd15052105125910591ef
  url: https://release-changing.io/2
  version: 2.0.0
- name: release-new
  sha1: d5435463432ea9aec93dd15052105125910591ef
  url: https://release-new.io/3
  version: 3.0.0
`))
		})
	})
})
