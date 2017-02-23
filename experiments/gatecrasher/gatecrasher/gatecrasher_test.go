package gatecrasher_test

import (
	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher"
	"gopkg.in/jarcoal/httpmock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type loggerMock struct {
	Messages []interface{}
}

func (l *loggerMock) Printf(format string, v ...interface{}) {
	l.Messages = append(l.Messages, []interface{}{format, v})
}

var ourLogger *loggerMock
var _ = Describe("Gatecrasher", func() {
	BeforeEach(func() {
		ourLogger = new(loggerMock)
	})
	Context("when the endpoint is good", func() {
		It("makes an https request", func() {
			url := "https://api.example.com/v2/info"
			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(200, `[{"id": 1, "name": "ALL YOUR INFOS"}]`))

			resp := gatecrasher.Run(url, ourLogger)
			Expect(len(ourLogger.Messages)).To(Equal(1))
			Expect(resp).To(Equal(200))
		})
		It("logs the request", func() {
			url := "https://api.example.com/v2/info"
			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(200, `[{"id": 1, "name": "ALL YOUR INFOS"}]`))

			resp := gatecrasher.Run(url, ourLogger)
			Expect(len(ourLogger.Messages)).To(Equal(1))
			Expect(resp).To(Equal(200))
		})

	})
	Context("when the endpoint is bad", func() {
		It("makes an https request", func() {
			url := "https://api.example.com/v2/info"
			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(502, `[{"id": 1, "name": "ALL YOUR INFOS"}]`))

			resp := gatecrasher.Run(url, ourLogger)
			Expect(len(ourLogger.Messages)).To(Equal(1))
			Expect(resp).To(Equal(502))
		})
	})
})
