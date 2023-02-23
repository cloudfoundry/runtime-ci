package concourseio

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConcourseio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Concourseio Suite")
}
