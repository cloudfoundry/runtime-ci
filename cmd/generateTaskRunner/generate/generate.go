package generate

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Package struct {
	name string
	dir  string

	inputs  []string
	outputs []string
}

func (p *Package) UnmarshalYAML(value *yaml.Node) error {
	var alias struct {
		Inputs  []map[string]string
		Outputs []map[string]string
	}
	err := value.Decode(&alias)
	if err != nil {
		return err
	}

	for _, input := range alias.Inputs {
		p.inputs = append(p.inputs, input["name"])
	}

	for _, output := range alias.Outputs {
		p.outputs = append(p.outputs, output["name"])
	}

	return nil
}

func NewPackage(outputPackageName, taskYAMLPath string) (Package, error) {
	pkg := Package{
		name: outputPackageName,
		dir:  outputPackageName,
	}

	contents, err := ioutil.ReadFile(taskYAMLPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(contents, &pkg)
	if err != nil {
		panic(err)
	}

	return pkg, nil
}

func (p *Package) SetOutputDir(outputDir string) {
	p.dir = outputDir
}

func (p Package) Write() error {
	err := os.MkdirAll(p.dir, 0755)
	if err != nil {
		panic(err)
	}

	return nil
}
