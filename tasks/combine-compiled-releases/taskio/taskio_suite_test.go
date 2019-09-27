package taskio

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTaskio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Taskio Suite")
}
