package manifest

import yaml "gopkg.in/yaml.v2"

func SetYAMLMarshal(f func(interface{}) ([]byte, error)) {
	YamlMarshal = f
}

func ResetYAMLMarshal() {
	YamlMarshal = yaml.Marshal
}
