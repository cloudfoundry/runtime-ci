package bosh

type Release struct {
	Name     string
	SHA1     string
	Stemcell Stemcell `yaml:",omitempty"`
	URL      string
	Version  string
}
