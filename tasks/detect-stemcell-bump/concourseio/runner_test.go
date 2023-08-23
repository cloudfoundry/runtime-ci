package concourseio

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Runner", func() {
	var (
		buildDir string
	)

	BeforeEach(func() {
		var err error
		buildDir, err = os.MkdirTemp("", "concourseio-rootdir-")
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
				expectedCFDeploymentDir string
				expectedStemcellDir     string
				expectedBumpTypeDir     string
			)

			BeforeEach(func() {
				expectedCFDeploymentDir = filepath.Join(buildDir, "cf-deployment")
				Expect(os.Mkdir(expectedCFDeploymentDir, 0777)).To(Succeed())
				expectedStemcellDir = filepath.Join(buildDir, "stemcell")
				Expect(os.Mkdir(expectedStemcellDir, 0777)).To(Succeed())
				expectedBumpTypeDir = filepath.Join(buildDir, "stemcell-bump-type")
				Expect(os.Mkdir(expectedBumpTypeDir, 0777)).To(Succeed())

			})

			It("will instantiate the runner", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualRunner).To(Equal(Runner{
					In: Inputs{
						cfDeploymentDir: expectedCFDeploymentDir,
						stemcellDir:     expectedStemcellDir,
					},
					Out: Outputs{
						bumpTypeDir: expectedBumpTypeDir,
					},
				}))
			})
		})

		Context("when the cf-deployment input is missing", func() {
			It("will return an error", func() {
				Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'cf-deployment' in build directory '%s'", buildDir)))
			})
		})

		Context("when the stemcell input is missing", func() {
			var (
				cfDeploymentDir string
			)

			BeforeEach(func() {
				cfDeploymentDir = filepath.Join(buildDir, "cf-deployment")
				Expect(os.Mkdir(cfDeploymentDir, 0777)).To(Succeed())
			})

			It("will return an error", func() {
				Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'stemcell' in build directory '%s'", buildDir)))
			})
		})

		Context("when the stemcell-bump-type output is missing", func() {
			var (
				cfDeploymentDir string
				stemcellDir     string
			)

			BeforeEach(func() {
				cfDeploymentDir = filepath.Join(buildDir, "cf-deployment")
				Expect(os.Mkdir(cfDeploymentDir, 0777)).To(Succeed())
				stemcellDir = filepath.Join(buildDir, "stemcell")
				Expect(os.Mkdir(stemcellDir, 0777)).To(Succeed())
			})

			It("will return an error", func() {
				Expect(actualErr).To(MatchError(fmt.Sprintf("missing sub directory 'stemcell-bump-type' in build directory '%s'", buildDir)))
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
				Expect(os.WriteFile(filepath.Join(stemcellDir, "version"), []byte("some-version"), 0777)).
					To(Succeed())
				Expect(os.WriteFile(filepath.Join(stemcellDir, "url"), []byte("https://s3.amazonaws.com/some-stemcell/stuff-ubuntu-some-os-go_agent.tgz"), 0777)).
					To(Succeed())
			})

			It("sets the stemcell OS and version", func() {
				Expect(actualErr).ToNot(HaveOccurred())

				Expect(runner.stemcell).To(Equal(bosh.Stemcell{OS: "ubuntu-some-os", Version: "some-version"}))
			})
		})

		Context("when the stemcell dir is missing the version file", func() {
			It("returns a missing files err", func() {
				Expect(actualErr).To(MatchError("missing files: 'stemcell/version'"))
			})
		})

		Context("when the stemcell dir is missing the url file", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(filepath.Join(stemcellDir, "version"), []byte("some-version"), 0777)).
					To(Succeed())
			})

			It("returns a missing files err", func() {
				Expect(actualErr).To(MatchError("missing files: 'stemcell/url'"))
			})
		})

		Context("when the stemcell dir contains all necessary files, but refers to an unsupported stemcell", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(filepath.Join(stemcellDir, "version"), []byte("some-version"), 0777)).
					To(Succeed())
				Expect(os.WriteFile(filepath.Join(stemcellDir, "url"), []byte("https://s3.amazonaws.com/some-stemcell/stuff-windows-some-os-go_agent.tgz"), 0777)).
					To(Succeed())
			})

			It("returns an unsupported stemcell type error", func() {
				Expect(actualErr).To(MatchError("stemcell URL does not contain a supported os (i.e. ubuntu): https://s3.amazonaws.com/some-stemcell/stuff-windows-some-os-go_agent.tgz"))
			})
		})
	})

	Describe("ReadCFDeploymentStemcell", func() {
		var (
			runner    Runner
			actualErr error
		)

		JustBeforeEach(func() {
			actualErr = runner.ReadCFDeploymentStemcell()
		})

		Context("when there is a stemcell defined in the manifest", func() {
			BeforeEach(func() {
				cfDeploymentDir := filepath.Join(buildDir, "cf-deployment")
				Expect(os.Mkdir(cfDeploymentDir, 0777)).
					To(Succeed())

				runner = Runner{In: Inputs{cfDeploymentDir: cfDeploymentDir}}
				cfDeploymentManifest := bosh.Manifest{
					Stemcells: []bosh.Stemcell{
						{
							OS:      "ubuntu-some-os",
							Version: "some-version-in-manifest",
						},
					},
				}
				cfDeploymentManifestContent, err := yaml.Marshal(cfDeploymentManifest)
				Expect(err).ToNot(HaveOccurred())

				err = os.WriteFile(filepath.Join(cfDeploymentDir, "cf-deployment.yml"), cfDeploymentManifestContent, 0644)
				Expect(err).ToNot(HaveOccurred())
			})

			It("sets the stemcell OS and version", func() {
				Expect(actualErr).ToNot(HaveOccurred())

				Expect(runner.manifestStemcell).To(Equal(bosh.Stemcell{OS: "ubuntu-some-os", Version: "some-version-in-manifest"}))
			})
		})

		Context("When the cf-deployment.yml file is missing", func() {
			It("returns a missing files err", func() {
				Expect(actualErr).To(MatchError("missing files: 'cf-deployment/cf-deployment.yml'"))
			})
		})

		Context("When the cf-deployment.yml file is malformed", func() {
			BeforeEach(func() {
				cfDeploymentDir := filepath.Join(buildDir, "cf-deployment")
				Expect(os.Mkdir(cfDeploymentDir, 0777)).
					To(Succeed())

				runner = Runner{In: Inputs{cfDeploymentDir: cfDeploymentDir}}

				err := os.WriteFile(filepath.Join(cfDeploymentDir, "cf-deployment.yml"), []byte("%%%"), 0644)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns a missing files err", func() {
				Expect(actualErr).To(MatchError("yaml: could not find expected directive name"))
			})
		})
	})

	Describe("DetectStemcellBump", func() {
		var (
			runner    Runner
			actualErr error
		)

		JustBeforeEach(func() {
			actualErr = runner.DetectStemcellBump()
		})

		Context("when the new stemcell is a forward bump", func() {
			BeforeEach(func() {
				runner.manifestStemcell = bosh.Stemcell{OS: "some-ubuntu", Version: "456.40"}
				runner.stemcell = bosh.Stemcell{OS: "some-ubuntu", Version: "457.0"}
			})

			It("returns the appropriate bump type", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(runner.bumpType).To(Equal("major"))
			})
		})

		Context("when the new stemcell is NOT a forward bump", func() {
			BeforeEach(func() {
				runner.manifestStemcell = bosh.Stemcell{OS: "some-ubuntu", Version: "500.0"}
				runner.stemcell = bosh.Stemcell{OS: "some-ubuntu", Version: "400.0"}
			})

			It("returns a non-forward bump error", func() {
				Expect(actualErr).To(HaveOccurred())
			})
		})
	})

	Describe("WriteStemcellBumpTypeToFile", func() {
		var (
			runner              Runner
			actualErr           error
			expectedBumpTypeDir string
		)

		JustBeforeEach(func() {
			actualErr = runner.WriteStemcellBumpTypeToFile()
		})

		BeforeEach(func() {
			expectedBumpTypeDir = filepath.Join(buildDir, "stemcell-bump-type")
			Expect(os.Mkdir(expectedBumpTypeDir, 0777)).To(Succeed())
			runner.Out = Outputs{expectedBumpTypeDir}
		})

		It("writes the bump type to the `bump-type` file", func() {
			runner.bumpType = "major"
			actualErr = runner.WriteStemcellBumpTypeToFile()
			Expect(actualErr).ToNot(HaveOccurred())

			actualResultContent, err := os.ReadFile(filepath.Join(expectedBumpTypeDir, "result"))
			Expect(err).ToNot(HaveOccurred())

			Expect(string(actualResultContent)).To(Equal("major"))
		})
	})
})
