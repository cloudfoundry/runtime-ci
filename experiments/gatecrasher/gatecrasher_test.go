package gatecrasher_test

import (
	"net/http"

	"github.com/krishicks/runtime-ci/experiments/gatecrasher"
	"github.com/krishicks/runtime-ci/experiments/gatecrasher/gatecrasherfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Gatecrasher", func() {
	var (
		fakeLogger *gatecrasherfakes.FakeLogger
		url        string
		server     *ghttp.Server
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		url = server.URL()
		fakeLogger = &gatecrasherfakes.FakeLogger{}

		server.AppendHandlers(
			ghttp.VerifyRequest("GET", "/"),
			ghttp.RespondWith(http.StatusOK, ""),
		)
	})

	AfterEach(func() {
		server.Close()
	})

	It("makes an https request", func() {
		gatecrasher.Run(url, fakeLogger)
		Expect(server.ReceivedRequests()).To(HaveLen(1))
	})

	It("logs the request", func() {
		gatecrasher.Run(url, fakeLogger)
		Expect(fakeLogger.PrintfCallCount()).To(Equal(1))
	})

	It("returns whatever response code it got", func() {
		statusCode := gatecrasher.Run(url, fakeLogger)
		Expect(statusCode).To(Equal(http.StatusOK))
	})
})
