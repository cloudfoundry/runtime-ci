package common

type MarshalFunc func(interface{}) ([]byte, error)
type UnmarshalFunc func([]byte, interface{}) error

type Release struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Version string `yaml:"version"`
	SHA1    string `yaml:"sha1"`
}