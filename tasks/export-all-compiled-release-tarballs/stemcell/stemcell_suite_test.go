package stemcell_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStemcells(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stemcells Suite")
}
