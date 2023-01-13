package compiledreleasesops_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCompiledreleasesops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Compiledreleasesops Suite")
}
