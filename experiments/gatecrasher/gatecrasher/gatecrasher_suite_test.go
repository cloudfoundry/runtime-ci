package gatecrasher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/jarcoal/httpmock.v1"

	"testing"
)

var _ = BeforeSuite(func() {
	// block all HTTP requests
	httpmock.Activate()
})

var _ = BeforeEach(func() {
	// remove any mocks
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

func TestGatecrasher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gatecrasher Suite")
}
