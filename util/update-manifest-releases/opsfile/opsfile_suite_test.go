package opsfile_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOpsFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "opsfile")
}
