package gatecrasher_test

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/config"
	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/gatecrasher/gatecrasherfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func addGoodHandlers(number int, server *ghttp.Server) {
	for i := 0; i < number; i++ {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/v2/info"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, ""),
			))
	}
}

var _ = Describe("Gatecrasher", func() {
	var fakeServer *ghttp.Server
	var fakeLogger *gatecrasherfakes.FakeLogger
	var configStruct config.Config
	var goodUrl string
	var badUrl string

	Describe("running for the configured number of requests", func() {
		BeforeEach(func() {
			fakeLogger = new(gatecrasherfakes.FakeLogger)
			fakeServer = ghttp.NewTLSServer()
			addGoodHandlers(2, fakeServer)
			configStruct = config.Load()
			configStruct.Target = fakeServer.URL() + "/v2/info"
		})

		Context("when the configuration specifies total run count", func() {
			BeforeEach(func() {
				configStruct.Total_number_of_requests = 2
			})
			It("runs the specified number of times", func() {
				gatecrasher.Run(configStruct, fakeLogger)
				Expect(fakeServer.ReceivedRequests()).To(HaveLen(2))
			})
		})

		AfterEach(func() {
			fakeServer.Close()
		})

		PContext("when the configuration a negative run count", func() {
			It("runs many times, maybe forever", func() {
				// We don't actually know how to test this well right now
				// It does, though.
			})
		})
	})

	Describe("logging 200s and 502s", func() {

		BeforeEach(func() {
			fakeLogger = new(gatecrasherfakes.FakeLogger)
			fakeServer = ghttp.NewTLSServer()
			goodUrl = fakeServer.URL() + "/v2/info"
			badUrl = fakeServer.URL() + "/v2/bad-info"
			configStruct = config.Load()
			configStruct.Total_number_of_requests = 1
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
				configStruct.Target = goodUrl
			})
			It("makes an https request", func() {
				gatecrasher.Run(configStruct, fakeLogger)
				Expect(fakeServer.ReceivedRequests()).Should(HaveLen(1))
			})

			It("logs the request in the correct format with necessary info", func() {
				loggedEvent := gatecrasher.EventLog{}
				gatecrasher.Run(configStruct, fakeLogger)
				format, args := fakeLogger.PrintfArgsForCall(0)
				Expect(format).To(Equal("%s"))
				Expect(len(args)).To(Equal(1))

				json.Unmarshal(args[0].([]byte), &loggedEvent)
				Expect(loggedEvent.URL).To(Equal(goodUrl))
				Expect(loggedEvent.StatusCode).To(Equal(http.StatusOK))
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
				configStruct.Target = badUrl
			})

			It("makes an https request", func() {
				gatecrasher.Run(configStruct, fakeLogger)
				Expect(fakeServer.ReceivedRequests()).Should(HaveLen(1))
			})

			It("logs the request in the correct format with necessary info", func() {
				loggedEvent := gatecrasher.EventLog{}
				gatecrasher.Run(configStruct, fakeLogger)
				format, args := fakeLogger.PrintfArgsForCall(0)
				Expect(format).To(Equal("%s"))
				Expect(len(args)).To(Equal(1))

				json.Unmarshal(args[0].([]byte), &loggedEvent)
				Expect(loggedEvent.URL).To(Equal(badUrl))
				Expect(loggedEvent.StatusCode).To(Equal(http.StatusBadGateway))
			})
		})
	})
	Describe("request summaries", func() {
		Context("when summaries are configured", func() {
			BeforeEach(func() {
				configStruct.Total_number_of_requests = 1
			})

			It("Logs a summary every configured interval", func() {
			})
		})
	})
})
