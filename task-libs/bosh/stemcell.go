package bosh

type Stemcell struct {
	OS      string
	Version string
}

// NewStemcellFromInput creates a Stemcell from a stemcell concourse resource
func NewStemcellFromInput(stemcellDir string) (Stemcell, error) {
	return Stemcell{}, nil
}
