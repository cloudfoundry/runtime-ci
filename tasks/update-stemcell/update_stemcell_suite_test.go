package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUpdateStemcell(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateStemcell Suite")
}
