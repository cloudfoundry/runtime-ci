package concourseio

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"

	. "github.com/onsi/ginkgo/v2"
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
			actualRunner Runner
			actualErr    error
		)

		JustBeforeEach(func() {
			actualRunner, actualErr = NewRunner(buildDir)
		})

		Context("when all directories exist", func() {
			var (
				expectedCFDeploymentDir        string
				expectedStemcellDir            string
				expectedUpdatedCFDeploymentDir string
			)

			BeforeEach(func() {
				expectedCFDeploymentDir = filepath.Join(buildDir, "cf-deployment")
				Expect(os.Mkdir(expectedCFDeploymentDir, 0777)).To(Succeed())
				expectedStemcellDir = filepath.Join(buildDir, "stemcell")
				Expect(os.Mkdir(expectedStemcellDir, 0777)).To(Succeed())
				expectedUpdatedCFDeploymentDir = filepath.Join(buildDir, "updated-cf-deployment")
				Expect(os.Mkdir(expectedUpdatedCFDeploymentDir, 0777)).To(Succeed())
			})

			It("will instantiate the runner", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualRunner).To(Equal(Runner{
					In: Inputs{
						cfDeploymentDir: expectedCFDeploymentDir,
						stemcellDir:     expectedStemcellDir,
					},
					Out: Outputs{
						UpdatedCFDeploymentDir: expectedUpdatedCFDeploymentDir,
					},
				}))
			})
		})

		Context("when some directories are missing", func() {
			var (
				expectedUpdatedCFDeploymentDir string
			)

			Context("when cf-deployment dir is missing", func() {
				BeforeEach(func() {
					expectedUpdatedCFDeploymentDir = filepath.Join(buildDir, "updated-cf-deployment")
					Expect(os.Mkdir(expectedUpdatedCFDeploymentDir, 0777)).To(Succeed())
				})

				It("will fail stating all the missing directories", func() {
					Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'cf-deployment' in build directory '%s'", buildDir)))
				})
			})
		})
	})

	Describe("ReadStemcell", func() {
		var (
			runner      Runner
			stemcellDir string

			actualErr error
		)

		BeforeEach(func() {
			stemcellDir = filepath.Join(buildDir, "stemcell")
			Expect(os.Mkdir(stemcellDir, 0777)).To(Succeed())

			runner = Runner{In: Inputs{stemcellDir: stemcellDir}}
		})

		JustBeforeEach(func() {
			actualErr = runner.ReadStemcell()
		})

		Context("when the stemcell dir contains all necessary files", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(filepath.Join(stemcellDir, "version"), []byte("some-version"), 0777)).
					To(Succeed())
				Expect(ioutil.WriteFile(filepath.Join(stemcellDir, "url"), []byte("https://s3.amazonaws.com/some-stemcell/stuff-ubuntu-some-os-go_agent.tgz"), 0777)).
					To(Succeed())
			})

			It("sets the stemcell OS and Version", func() {
				Expect(actualErr).ToNot(HaveOccurred())

				Expect(runner.stemcell).To(Equal(bosh.Stemcell{OS: "ubuntu-some-os", Version: "some-version"}))
			})
		})

		Context("when the stemcell dir is missing some files", func() {
			It("returns a missing files err", func() {
				Expect(actualErr).To(MatchError("missing files: 'stemcell/version'"))
			})
		})
	})

	Describe("UpdateManifest", func() {
		var (
			runner Runner

			expectedCFDeploymentDir        string
			expectedUpdatedCFDeploymentDir string
			expectedStemcell               bosh.Stemcell

			manifestUpdateSpy        UpdateFunc
			manifestUpdateFileOutput []byte

			actualInFile   []byte
			actualStemcell bosh.Stemcell

			actualErr error
		)

		BeforeEach(func() {
			expectedCFDeploymentDir = filepath.Join(buildDir, "cf-deployment")
			Expect(os.Mkdir(expectedCFDeploymentDir, 0777)).To(Succeed())
			expectedUpdatedCFDeploymentDir = filepath.Join(buildDir, "updated-cf-deployment")
			Expect(os.Mkdir(expectedUpdatedCFDeploymentDir, 0777)).To(Succeed())

			expectedStemcell = bosh.Stemcell{OS: "gundam", Version: "1.1.0"}
			runner = Runner{
				stemcell: expectedStemcell,
				In:       Inputs{cfDeploymentDir: expectedCFDeploymentDir},
				Out:      Outputs{UpdatedCFDeploymentDir: expectedUpdatedCFDeploymentDir},
			}

			manifestUpdateSpy = func(file []byte, stemcell bosh.Stemcell) ([]byte, error) {
				actualInFile = file
				actualStemcell = stemcell

				return manifestUpdateFileOutput, nil
			}
		})

		JustBeforeEach(func() {
			actualErr = runner.UpdateManifest(manifestUpdateSpy)
		})

		Context("When there exists a valid manifest", func() {
			var (
				expectedInFile []byte
			)

			BeforeEach(func() {
				expectedInFile = []byte("This is my manifest")
				manifestPath := filepath.Join(expectedCFDeploymentDir, "cf-deployment.yml")
				Expect(ioutil.WriteFile(manifestPath, expectedInFile, 0777)).To(Succeed())
			})

			It("Updates the manifest", func() {
				Expect(actualStemcell).To(Equal(expectedStemcell))
				Expect(actualInFile).To(Equal(expectedInFile))
			})

			Context("when the manifest update function returns an updated Manifest", func() {
				BeforeEach(func() {
					manifestUpdateFileOutput = []byte("updated manifest")
					expectedInFile = []byte("This is the manifest in my updatedCFDeploymentDir")

					manifestPath := filepath.Join(expectedUpdatedCFDeploymentDir, "cf-deployment.yml")
					Expect(ioutil.WriteFile(manifestPath, expectedInFile, 0777)).To(Succeed())
				})

				It("writes the file to the output file", func() {
					Expect(actualErr).ToNot(HaveOccurred())
					actualOutFile, err := ioutil.ReadFile(filepath.Join(expectedUpdatedCFDeploymentDir, "cf-deployment.yml"))
					Expect(err).ToNot(HaveOccurred())

					Expect(actualOutFile).To(Equal(manifestUpdateFileOutput))
				})
			})
		})
	})

	Describe("WriteCommitMessage", func() {
		var (
			runner            Runner
			expectedStemcell  bosh.Stemcell
			actualErr         error
			commitMessagePath string
		)

		BeforeEach(func() {
			commitMessagePath = filepath.Join(buildDir, "commit-message.txt")
			expectedStemcell = bosh.Stemcell{OS: "gundam", Version: "1.1.0"}
			runner = Runner{
				stemcell: expectedStemcell,
			}
		})

		JustBeforeEach(func() {
			actualErr = runner.WriteCommitMessage(commitMessagePath)
		})

		It("Writes a message with the new stemcell version", func() {
			Expect(actualErr).ToNot(HaveOccurred())

			actualCommitMessage, err := ioutil.ReadFile(commitMessagePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(actualCommitMessage)).To(Equal(fmt.Sprintf("Update stemcell to %s \"%s\"",
				expectedStemcell.OS, expectedStemcell.Version)))
		})
	})
})
