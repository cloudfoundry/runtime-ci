package compiledrelease_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCompiledrelease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Compiledrelease Suite")
}
