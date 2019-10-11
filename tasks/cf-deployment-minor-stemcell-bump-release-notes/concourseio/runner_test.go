package concourseio_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	"github.com/cloudfoundry/runtime-ci/tasks/cf-deployment-minor-stemcell-bump-release-notes/concourseio"
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
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(buildDir)).To(Succeed())
	})

	Describe("NewRunner", func() {
		var (
			expectedCFDeploymentDir   string
			expectedReleaseVersionDir string
			expectedStemcellDir       string
			expectedReleaseNotesDir   string

			actualRunner concourseio.Runner
			actualErr    error
		)

		JustBeforeEach(func() {
			actualRunner, actualErr = concourseio.NewRunner(buildDir)
		})

		Context("when all directories exist", func() {
			BeforeEach(func() {
				expectedCFDeploymentDir = filepath.Join(buildDir, "cf-deployment-master")
				Expect(os.Mkdir(expectedCFDeploymentDir, 0777)).To(Succeed())
				expectedReleaseVersionDir = filepath.Join(buildDir, "release-version")
				Expect(os.Mkdir(expectedReleaseVersionDir, 0777)).To(Succeed())
				expectedStemcellDir = filepath.Join(buildDir, "stemcell")
				Expect(os.Mkdir(expectedStemcellDir, 0777)).To(Succeed())
				expectedReleaseNotesDir = filepath.Join(buildDir, "cf-deployment-minor-stemcell-bump-release-notes")
				Expect(os.Mkdir(expectedReleaseNotesDir, 0777)).To(Succeed())
			})

			It("will instantiate the runner", func() {
				Expect(actualErr).NotTo(HaveOccurred())
				Expect(actualRunner).To(Equal(concourseio.Runner{
					In: concourseio.Inputs{
						CFDeploymentDir:   expectedCFDeploymentDir,
						ReleaseVersionDir: expectedReleaseVersionDir,
						StemcellDir:       expectedStemcellDir,
					},
					Out: concourseio.Outputs{
						ReleaseNotesDir: expectedReleaseNotesDir,
					},
				}))
			})
		})

		Context("when some directories are missing", func() {
			Context("when cf-deployment-master dir is missing", func() {
				BeforeEach(func() {
					expectedStemcellDir = filepath.Join(buildDir, "stemcell")
					Expect(os.Mkdir(expectedStemcellDir, 0777)).To(Succeed())
					expectedReleaseNotesDir = filepath.Join(buildDir, "cf-deployment-minor-stemcell-bump-release-notes")
					Expect(os.Mkdir(expectedReleaseNotesDir, 0777)).To(Succeed())
				})

				It("will fail stating all the missing directories", func() {
					Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'cf-deployment-master' in build directory '%s'", buildDir)))
				})
			})

			Context("when stemcell dir is missing", func() {
				BeforeEach(func() {
					expectedCFDeploymentDir = filepath.Join(buildDir, "cf-deployment-master")
					Expect(os.Mkdir(expectedCFDeploymentDir, 0777)).To(Succeed())
					expectedReleaseVersionDir = filepath.Join(buildDir, "release-version")
					Expect(os.Mkdir(expectedReleaseVersionDir, 0777)).To(Succeed())
					expectedStemcellDir = filepath.Join(buildDir, "stemcell")
					expectedReleaseNotesDir = filepath.Join(buildDir, "cf-deployment-minor-stemcell-bump-release-notes")
					Expect(os.Mkdir(expectedReleaseNotesDir, 0777)).To(Succeed())
				})

				It("will fail stating all the missing directories", func() {
					Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'stemcell' in build directory '%s'", buildDir)))
				})
			})
		})
	})

	Describe("ReadStemcellInfoFromManifest", func() {
		var (
			actualStemcell  bosh.Stemcell
			actualErr       error
			cfDeploymentDir string
			stemcellAlias   string
			runner          concourseio.Runner
		)

		BeforeEach(func() {
			cfDeploymentDir = filepath.Join(buildDir, "cf-deployment-master")
			Expect(os.Mkdir(cfDeploymentDir, 0777)).To(Succeed())

			runner = concourseio.Runner{
				In: concourseio.Inputs{
					CFDeploymentDir: cfDeploymentDir,
				},
			}

			stemcellAlias = "default"
		})

		JustBeforeEach(func() {
			actualStemcell, actualErr = runner.ReadStemcellInfoFromManifest(stemcellAlias)
		})

		Context("happy path", func() {
			BeforeEach(func() {
				cfDeploymentManifest := []byte(`---
stemcells:
- alias: default
  os: some-ubuntu
  version: some-version
- alias: other-stemcell
  os: some-other-os
  version: some-other-version
`)
				Expect(ioutil.WriteFile(filepath.Join(cfDeploymentDir, "cf-deployment.yml"), cfDeploymentManifest, 0644)).To(Succeed())
			})

			It("returns a stemcell struct from the specified stemcell", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualStemcell).To(Equal(bosh.Stemcell{
					Alias:   "default",
					OS:      "some-ubuntu",
					Version: "some-version",
				}))
			})
		})

		Context("failure cases", func() {
			Context("when the manifest does not exist", func() {
				It("returns a wrapped error", func() {
					Expect(actualErr).To(HaveOccurred())
					Expect(actualErr.Error()).To(ContainSubstring("failed to read cf-deployment.yml"))
					innerErr := errors.Unwrap(actualErr)
					Expect(innerErr.Error()).To(ContainSubstring("no such file or directory"))
				})
			})

			Context("when the manifest yml is invalid", func() {
				BeforeEach(func() {
					cfDeploymentManifest := []byte(`%%%`)
					Expect(ioutil.WriteFile(filepath.Join(cfDeploymentDir, "cf-deployment.yml"), cfDeploymentManifest, 0644)).To(Succeed())
				})

				It("returns a wrapped yaml error", func() {
					Expect(actualErr).To(HaveOccurred())
					Expect(actualErr.Error()).To(ContainSubstring("failed to unmarshal cf-deployment.yml"))
					innerErr := errors.Unwrap(actualErr)
					Expect(innerErr.Error()).To(ContainSubstring("yaml: could not find expected directive name"))
				})
			})

			Context("when the stemcell alias is not found in the manifest", func() {
				BeforeEach(func() {
					cfDeploymentManifest := []byte(`---
stemcells:
- alias: default
  os: some-ubuntu
  version: some-version
`)
					Expect(ioutil.WriteFile(filepath.Join(cfDeploymentDir, "cf-deployment.yml"), cfDeploymentManifest, 0644)).To(Succeed())

					stemcellAlias = "missing-alias"
				})

				It("returns an error", func() {
					Expect(actualErr).To(HaveOccurred())
					Expect(actualErr.Error()).To(ContainSubstring(`failed to find stemcell version for alias "missing-alias"`))
				})
			})
		})
	})

	Describe("ValidateStemcellBump", func() {
		var (
			oldStemcell bosh.Stemcell
			newStemcell bosh.Stemcell
			actualErr   error
		)

		BeforeEach(func() {
			oldStemcell = bosh.Stemcell{
				Alias:   "default",
				OS:      "some-ubuntu-os",
				Version: "some-old-version",
			}
		})

		JustBeforeEach(func() {
			runner := concourseio.Runner{}
			actualErr = runner.ValidateStemcellBump(oldStemcell, newStemcell)
		})

		Context("valid bump", func() {
			BeforeEach(func() {
				newStemcell = bosh.Stemcell{
					Alias:   "default",
					OS:      "some-ubuntu-os",
					Version: "some-new-version",
				}
			})

			It("does not return an error", func() {
				Expect(actualErr).NotTo(HaveOccurred())
			})
		})

		Context("os mismatch", func() {
			BeforeEach(func() {
				newStemcell = bosh.Stemcell{
					OS:      "some-windows-os",
					Version: "some-version",
				}
			})

			It("returns an error", func() {
				Expect(actualErr).To(MatchError(`stemcell os mismatch: found "some-ubuntu-os" and "some-windows-os"`))
			})
		})
	})

	Describe("GenerateReleaseNotes", func() {
		var (
			oldStemcell             bosh.Stemcell
			newStemcell             bosh.Stemcell
			expectedReleaseNotesDir string
			actualErr               error
		)

		BeforeEach(func() {
			oldStemcell = bosh.Stemcell{
				Alias:   "default",
				OS:      "ubuntu-xenial",
				Version: "some-old-version",
			}
		})

		JustBeforeEach(func() {
			runner := concourseio.Runner{
				Out: concourseio.Outputs{
					ReleaseNotesDir: expectedReleaseNotesDir,
				},
			}
			actualErr = runner.GenerateReleaseNotes(oldStemcell, newStemcell)
		})

		Context("happy path", func() {
			BeforeEach(func() {
				newStemcell = bosh.Stemcell{
					Alias:   "default",
					OS:      "ubuntu-xenial",
					Version: "some-new-version",
				}

				expectedReleaseNotesDir = filepath.Join(buildDir, "cf-deployment-minor-stemcell-bump-release-notes")
				Expect(os.Mkdir(expectedReleaseNotesDir, 0777)).To(Succeed())
			})

			It("writes the release notes output to a file", func() {
				Expect(actualErr).NotTo(HaveOccurred())
				outputFilePath := filepath.Join(expectedReleaseNotesDir, "body.txt")
				Expect(outputFilePath).To(BeAnExistingFile())
				releaseNotesContent, err := ioutil.ReadFile(outputFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(releaseNotesContent)).To(Equal(`## Stemcell Updates
| Release | Old Version | New Version |
| - | - | - |
| ubuntu xenial | some-old-version | some-new-version |
`))
			})
		})

		Context("when the release notes file cannot be written", func() {
			BeforeEach(func() {
				expectedReleaseNotesDir = filepath.Join(buildDir, "missing-dir")
			})

			It("returns the error", func() {
				Expect(actualErr).To(HaveOccurred())
				Expect(actualErr.Error()).To(ContainSubstring("failed to write release notes file"))
				innerErr := errors.Unwrap(actualErr)
				Expect(innerErr.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})

	Describe("GenerateReleaseName", func() {
		var (
			expectedReleaseVersionDir string
			expectedReleaseNotesDir   string
			actualErr                 error
			version                   string
		)

		BeforeEach(func() {
			version = "1.2.3"
			expectedReleaseVersionDir = filepath.Join(buildDir, "release-version")
			Expect(os.Mkdir(expectedReleaseVersionDir, 0777)).To(Succeed())
			expectedReleaseNotesDir = filepath.Join(buildDir, "cf-deployment-minor-stemcell-bump-release-notes")
			Expect(os.Mkdir(expectedReleaseNotesDir, 0777)).To(Succeed())
		})

		JustBeforeEach(func() {
			runner := concourseio.Runner{
				In: concourseio.Inputs{
					ReleaseVersionDir: expectedReleaseVersionDir,
				},
				Out: concourseio.Outputs{
					ReleaseNotesDir: expectedReleaseNotesDir,
				},
			}
			actualErr = runner.GenerateReleaseName()
		})

		Context("happy path", func() {
			BeforeEach(func() {
				err := ioutil.WriteFile(filepath.Join(expectedReleaseVersionDir, "version"), []byte(version), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Writes the v-prefixed release version to the name file", func() {
				Expect(actualErr).NotTo(HaveOccurred())
				outputFilePath := filepath.Join(expectedReleaseNotesDir, "name.txt")
				Expect(outputFilePath).To(BeAnExistingFile())
				releaseNotesContent, err := ioutil.ReadFile(outputFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(releaseNotesContent)).To(Equal(fmt.Sprintf("v%s", version)))
			})
		})

		Context("when the release version file does not exist", func() {
			It("returns the error", func() {
				Expect(actualErr).To(HaveOccurred())
				Expect(actualErr.Error()).To(ContainSubstring("failed to read release version"))
				innerErr := errors.Unwrap(actualErr)
				Expect(innerErr.Error()).To(ContainSubstring("no such file or directory"))
			})
		})

		Context("when the release name file cannot be written", func() {
			BeforeEach(func() {
				err := ioutil.WriteFile(filepath.Join(expectedReleaseVersionDir, "version"), []byte(version), 0644)
				Expect(err).NotTo(HaveOccurred())

				expectedReleaseNotesDir = filepath.Join(buildDir, "missing-dir")
			})

			It("returns the error", func() {
				Expect(actualErr).To(HaveOccurred())
				Expect(actualErr.Error()).To(ContainSubstring("failed to write release name file"))
				innerErr := errors.Unwrap(actualErr)
				Expect(innerErr.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})
})
