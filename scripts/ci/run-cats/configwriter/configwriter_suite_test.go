package configwriter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConfigwriter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configwriter Suite")
}
