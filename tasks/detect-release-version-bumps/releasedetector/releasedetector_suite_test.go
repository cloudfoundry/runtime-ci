package releasedetector

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestReleasedetector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Releasedetector Suite")
}
