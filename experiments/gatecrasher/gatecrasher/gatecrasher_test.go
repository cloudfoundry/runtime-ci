package gatecrasher_test

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher/gatecrasherfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Gatecrasher", func() {
	var fakeServer *ghttp.Server
	var fakeLogger *gatecrasherfakes.FakeLogger
	var goodUrl string
	var badUrl string

	BeforeEach(func() {
		fakeLogger = new(gatecrasherfakes.FakeLogger)
		fakeServer = ghttp.NewTLSServer()
		goodUrl = fakeServer.URL() + "/v2/info"
		badUrl = fakeServer.URL() + "/v2/bad-info"
	})

	AfterEach(func() {
		fakeServer.Close()
	})

	Context("when the endpoint is good", func() {
		BeforeEach(func() {
			fakeServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/info"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, ""),
				),
			)
		})
		It("makes an https request", func() {
			resp := gatecrasher.Run(goodUrl, fakeLogger)

			Expect(fakeServer.ReceivedRequests()).Should(HaveLen(1))
			Expect(resp).To(Equal(http.StatusOK))
		})

		It("logs the request in the correct format with necessary info", func() {
			loggedEvent := gatecrasher.EventLog{}
			gatecrasher.Run(goodUrl, fakeLogger)
			format, args := fakeLogger.PrintfArgsForCall(0)
			Expect(format).To(Equal("%s"))
			Expect(len(args)).To(Equal(1))

			json.Unmarshal(args[0].([]byte), &loggedEvent)
			Expect(loggedEvent.URL).To(Equal(goodUrl))
			Expect(loggedEvent.StatusCode).To(Equal(http.StatusOK))

			flag := fakeLogger.SetFlagsArgsForCall(0)
			Expect(flag).To(Equal(log.LstdFlags))
		})
	})

	Context("when the endpoint is bad", func() {
		BeforeEach(func() {
			fakeServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/bad-info"),
					ghttp.RespondWith(http.StatusBadGateway, ""),
				),
			)
		})

		It("makes an https request", func() {
			resp := gatecrasher.Run(badUrl, fakeLogger)
			Expect(resp).To(Equal(http.StatusBadGateway))
		})
	})
})
