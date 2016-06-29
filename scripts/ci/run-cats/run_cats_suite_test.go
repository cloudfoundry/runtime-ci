package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRunCats(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RunCats Suite")
}
