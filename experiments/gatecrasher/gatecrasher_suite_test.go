package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGatecrasher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gatecrasher Suite")
}
