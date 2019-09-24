package concourseio

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"

	"github.com/cloudfoundry/runtime-ci/tasks/detect-release-version-bumps/concourseio/concourseiofakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	var (
		buildDir string
	)

	BeforeEach(func() {
		var err error
		buildDir, err = ioutil.TempDir("", "concourseio-rootdir-")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(buildDir)).To(Succeed())
	})

	Describe("NewRunner", func() {
		var (
			expectedCFDeploymentPrevDir string
			expectedCFDeploymentNextDir string
			expectedReleaseListDir      string

			actualRunner Runner
			actualErr    error
		)

		JustBeforeEach(func() {
			actualRunner, actualErr = NewRunner(buildDir)
		})

		Context("when all directories exist", func() {
			BeforeEach(func() {
				expectedCFDeploymentPrevDir = filepath.Join(buildDir, "cf-deployment-prev")
				Expect(os.Mkdir(expectedCFDeploymentPrevDir, 0777)).To(Succeed())

				expectedCFDeploymentNextDir = filepath.Join(buildDir, "cf-deployment-next")
				Expect(os.Mkdir(expectedCFDeploymentNextDir, 0777)).To(Succeed())

				expectedReleaseListDir = filepath.Join(buildDir, "release-list")
				Expect(os.Mkdir(expectedReleaseListDir, 0777)).To(Succeed())
			})

			It("will instantiate the runner", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualRunner).To(Equal(Runner{
					In: Inputs{
						cfDeploymentPrevDir: expectedCFDeploymentPrevDir,
						cfDeploymentNextDir: expectedCFDeploymentNextDir,
					},
					Out: Outputs{
						releaseListDir: expectedReleaseListDir,
					},
				}))
			})
		})

		Context("when some directories are missing", func() {
			Context("when cf-deployment-prev dir is missing", func() {
				BeforeEach(func() {
					expectedCFDeploymentNextDir = filepath.Join(buildDir, "cf-deployment-next")
					Expect(os.Mkdir(expectedCFDeploymentNextDir, 0777)).To(Succeed())

					expectedReleaseListDir = filepath.Join(buildDir, "release-list")
					Expect(os.Mkdir(expectedReleaseListDir, 0777)).To(Succeed())
				})

				It("will fail stating all the missing directories", func() {
					Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'cf-deployment-prev' in build directory '%s'", buildDir)))
				})
			})

			Context("when cf-deployment-next dir is missing", func() {
				BeforeEach(func() {
					expectedCFDeploymentPrevDir = filepath.Join(buildDir, "cf-deployment-prev")
					Expect(os.Mkdir(expectedCFDeploymentPrevDir, 0777)).To(Succeed())

					expectedReleaseListDir = filepath.Join(buildDir, "release-list")
					Expect(os.Mkdir(expectedReleaseListDir, 0777)).To(Succeed())
				})

				It("will fail stating all the missing directories", func() {
					Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'cf-deployment-next' in build directory '%s'", buildDir)))
				})
			})

			Context("when release-list dir is missing", func() {
				BeforeEach(func() {
					expectedCFDeploymentPrevDir = filepath.Join(buildDir, "cf-deployment-prev")
					Expect(os.Mkdir(expectedCFDeploymentPrevDir, 0777)).To(Succeed())

					expectedCFDeploymentNextDir = filepath.Join(buildDir, "cf-deployment-next")
					Expect(os.Mkdir(expectedCFDeploymentNextDir, 0777)).To(Succeed())
				})

				It("will fail stating all the missing directories", func() {
					Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'release-list' in build directory '%s'", buildDir)))
				})
			})
		})
	})

	Describe("Run", func() {
		var (
			detector       *concourseiofakes.FakeReleaseDetector
			runner         Runner
			releaseListDir string
			err            error
		)

		BeforeEach(func() {
			releaseListDir = filepath.Join(buildDir, "release-list")
			Expect(os.Mkdir(releaseListDir, 0777)).To(Succeed())
			detector = new(concourseiofakes.FakeReleaseDetector)
			runner = Runner{
				In:  Inputs{cfDeploymentPrevDir: "/prev-deployment", cfDeploymentNextDir: "/next-deployment"},
				Out: Outputs{releaseListDir: releaseListDir},
			}
		})

		JustBeforeEach(func() {
			err = runner.Run(detector)
		})

		Context("on success", func() {
			var (
				expectedContent []byte
			)

			BeforeEach(func() {
				expectedContent = []byte("this is a list of releases to be deployed")
				detector.WriteReturns(expectedContent, nil)
			})

			It("detects new/updated releases and writes a release list yaml file", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(detector.LoadCallCount()).To(Equal(1), "Expected to call Load")
				prevManifestPath, nextManifestPath := detector.LoadArgsForCall(0)
				Expect(prevManifestPath).To(Equal("/prev-deployment/cf-deployment.yml"))
				Expect(nextManifestPath).To(Equal("/next-deployment/cf-deployment.yml"))
				Expect(detector.DetectUpdatedReleasesCallCount()).To(Equal(1), "Expected to call DetectUpdatedReleases")
				Expect(detector.WriteCallCount()).To(Equal(1), "Expected to call Write")

				actualReleaseList, err := ioutil.ReadFile(filepath.Join(buildDir, "release-list", "releases.yml"))
				Expect(err).ToNot(HaveOccurred())

				Expect(actualReleaseList).To(Equal(expectedContent))
			})
		})

		Context("on failure", func() {
			Context("when load fails to read a file", func() {
				BeforeEach(func() {
					detector.LoadReturns(errors.New("failed to load"))
				})

				It("returns the error", func() {
					Expect(err).To(MatchError("failed to load"))
				})
			})

			Context("when write fails to generate the release list yaml", func() {
				BeforeEach(func() {
					detector.WriteReturns(nil, errors.New("failed to write"))
				})

				It("returns the error", func() {
					Expect(err).To(MatchError("failed to write"))
				})
			})

			Context("when the release list yaml file cannot be written", func() {
				BeforeEach(func() {
					runner = Runner{Out: Outputs{releaseListDir: "/does/not/exist"}}
				})
				It("returns the error", func() {
					Expect(err).To(MatchError("open /does/not/exist/releases.yml: no such file or directory"))
				})
			})
		})
	})
})
