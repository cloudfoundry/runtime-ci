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
			buildDir string
		)

		BeforeEach(func() {
			var err error
			buildDir, err = generateBuildDir(taskYAMLPath)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	JustBeforeEach(func() {
		runner, err = 
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
