package concourseio_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConcourseio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Concourseio Suite")
}
