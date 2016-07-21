package validationerrors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestValidationerrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validationerrors Suite")
}
