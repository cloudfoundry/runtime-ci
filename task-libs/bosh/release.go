package bosh

type Release struct {
	Name     string
	SHA1     string
	Stemcell Stemcell
	Version  string
	URL      string
}
