package taskio

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v2"
)

const taskYAMLPath = "../task.yml"

var _ = Describe("Runner", func() {
	Describe("NewRunner", func() {
		var (
			buildDir     string
			actualRunner Runner
			actualErr    error
		)

		BeforeEach(func() {
			var err error
			buildDir, err = generateBuildDir(taskYAMLPath)
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			actualRunner, actualErr = NewRunner(buildDir)
		})

		It("creates a new runner", func() {
			Expect(actualErr).ToNot(HaveOccurred())
			Expect(actualRunner).To(Equal(Runner{
				In: Inputs{
					compiledReleasesPrevDir: filepath.Join(buildDir, "compiled-releases-prev"),
					compiledReleasesNextDir: filepath.Join(buildDir, "compiled-releases-next"),
					cfDeploymentNextDir:     filepath.Join(buildDir, "cf-deployment-next"),
				},
				Out: Outputs{
					compliledReleasesCombinedDir: filepath.Join(buildDir, "compiled-releases-combined"),
				},
			}))
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			// create fake tarballcombiner
		})
	})
})

type task struct {
	Inputs  []input
	Outputs []output
}

type input struct {
	Name     string
	Optional bool
}

type output struct {
	Name string
}

func generateBuildDir(path string) (string, error) {
	tmp, err := ioutil.TempDir("", "build-dir-")
	if err != nil {
		return "", err
	}

	taskFile, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var taskYAML task
	err = yaml.Unmarshal(taskFile, &taskYAML)
	if err != nil {
		return "", err
	}

	for _, input := range taskYAML.Inputs {
		err := os.Mkdir(filepath.Join(tmp, input.Name), 0755)
		if err != nil {
			return "", err
		}
	}

	for _, output := range taskYAML.Outputs {
		err := os.Mkdir(filepath.Join(tmp, output.Name), 0755)
		if err != nil {
			return "", err
		}
	}

	return tmp, nil
}
