package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBBLMixin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "bbl-mixin-routing-release-lbs")
}
