package concourseio

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	Describe("NewRunner", func() {
		var (
			buildDir string

			actualRunner Runner
			actualErr    error
		)

		BeforeEach(func() {
			var err error
			buildDir, err = ioutil.TempDir("", "concourseio-rootdir-")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(buildDir)).To(Succeed())
		})

		JustBeforeEach(func() {
			// for key, val := range envs {
			// 	os.SetEnv(key, val)
			// }

			actualRunner, actualErr = NewRunner(buildDir)
		})

		Context("when all directories exist", func() {
			var (
				expectedCFDeploymentDir        string
				expectedCompiledReleasesDir    string
				expectedStemcellDir            string
				expectedUpdatedCFDeploymentDir string
			)

			BeforeEach(func() {
				expectedCFDeploymentDir = filepath.Join(buildDir, "cf-deployment")
				Expect(os.Mkdir(expectedCFDeploymentDir, 0777)).To(Succeed())
				expectedCompiledReleasesDir = filepath.Join(buildDir, "compiled-releases")
				Expect(os.Mkdir(expectedCompiledReleasesDir, 0777)).To(Succeed())
				expectedStemcellDir = filepath.Join(buildDir, "stemcell")
				Expect(os.Mkdir(expectedStemcellDir, 0777)).To(Succeed())
				expectedUpdatedCFDeploymentDir = filepath.Join(buildDir, "updated-cf-deployment")
				Expect(os.Mkdir(expectedUpdatedCFDeploymentDir, 0777)).To(Succeed())
			})

			It("will instantiate the runner", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualRunner).To(Equal(Runner{
					In: Inputs{
						cfDeploymentDir:     expectedCFDeploymentDir,
						compiledReleasesDir: expectedCompiledReleasesDir,
						stemcellDir:         expectedStemcellDir,
					},
					Out: Outputs{
						updatedCFDeploymentDir: expectedUpdatedCFDeploymentDir,
					},
				}))
			})
		})

		Context("when some directories are missing", func() {
			var (
				expectedCompiledReleasesDir    string
				expectedUpdatedCFDeploymentDir string
			)

			BeforeEach(func() {
				expectedCompiledReleasesDir = filepath.Join(buildDir, "compiled-releases")
				Expect(os.Mkdir(expectedCompiledReleasesDir, 0777)).To(Succeed())
				expectedUpdatedCFDeploymentDir = filepath.Join(buildDir, "updated-cf-deployment")
				Expect(os.Mkdir(expectedUpdatedCFDeploymentDir, 0777)).To(Succeed())
			})

			It("will fail stating all the missing directories", func() {
				Expect(actualErr).To(MatchError("missing directories: 'cf-deployment'"))
			})
		})
	})

	Describe("ReadStemcell", func() {
		var (
			runner      Runner
			buildDir    string
			stemcellDir string

			actualErr error
		)

		BeforeEach(func() {
			var err error
			buildDir, err = ioutil.TempDir("", "concourseio-stemcelldir-")
			Expect(err).ToNot(HaveOccurred())

			stemcellDir = filepath.Join(buildDir, "stemcell")
			Expect(os.Mkdir(stemcellDir, 0777)).To(Succeed())

			runner = Runner{In: Inputs{stemcellDir: stemcellDir}}
		})

		AfterEach(func() {
			Expect(os.RemoveAll(buildDir)).To(Succeed())
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

				Expect(runner.stemcell).To(Equal(manifest.Stemcell{OS: "ubuntu-some-os", Version: "some-version"}))
			})
		})

		Context("when the stemcell dir is missing some files", func() {
			It("returns a missing files err", func() {
				Expect(actualErr).To(MatchError("missing files: 'stemcell/version'"))
			})
		})
	})

	Describe("Update", func() {
		var (
			buildDir                string
			expectedCFDeploymentDir string
			runner                  Runner
			expectedStemcell        manifest.Stemcell
			expectedInFile          []byte

			manifestUpdateSpy UpdateFunc
			actualInFile      []byte
			actualStemcell    manifest.Stemcell
		)

		BeforeEach(func() {
			// MAKE A DIR with the cf-deployment yml
			var err error
			buildDir, err = ioutil.TempDir("", "concourseio-rootdir-")
			Expect(err).ToNot(HaveOccurred())
			expectedCFDeploymentDir = filepath.Join(buildDir, "cf-deployment")
			Expect(os.Mkdir(expectedCFDeploymentDir, 0777)).To(Succeed())
			expectedInFile = []byte("This is my manifest")
			manifestPath := filepath.Join(expectedCFDeploymentDir, "cf-deployment.yml")
			Expect(ioutil.WriteFile(manifestPath, expectedInFile, 0777)).
				To(Succeed())

			expectedStemcell = manifest.Stemcell{
				OS:      "gundam",
				Version: "1.1.0",
			}
			runner = Runner{
				stemcell: expectedStemcell,
				In: Inputs{
					cfDeploymentDir: expectedCFDeploymentDir,
				},
			}

			manifestUpdateSpy = func(file []byte, stemcell manifest.Stemcell) ([]byte, error) {
				actualInFile = file
				actualStemcell = stemcell

				return []byte("updated manifest"), nil
			}
		})

		AfterEach(func() {
			Expect(os.RemoveAll(buildDir)).To(Succeed())
		})

		Context("When ...", func() {
			It("Updates the manifest and compiled releases opsfiles", func() {
				Expect(runner.UpdateManifest(manifestUpdateSpy)).To(BeNil())

				Expect(actualStemcell).To(Equal(expectedStemcell))
				Expect(actualInFile).To(Equal(expectedInFile))

				//eventually write the new Thing to the output cf deployment file
			})

			It("writes the file to the output file", func() {

			})
		})
	})
})
