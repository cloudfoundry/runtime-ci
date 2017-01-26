package manifest

import yaml "gopkg.in/yaml.v2"

func SetYAMLMarshal(f func(interface{}) ([]byte, error)) {
	yamlMarshal = f
}

func ResetYAMLMarshal() {
	yamlMarshal = yaml.Marshal
}

func SetYAMLUnmarshal(f func([]byte, interface{}) error) {
	yamlUnmarshal = f
}

func ResetYAMLUnmarshal() {
	yamlUnmarshal = yaml.Unmarshal
}
