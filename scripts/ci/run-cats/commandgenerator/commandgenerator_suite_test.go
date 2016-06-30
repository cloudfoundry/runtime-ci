package commandgenerator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCommandgenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commandgenerator Suite")
}
