package config_test

import (
	"os"

	"github.com/cloudfoundry/runtime-ci/experiments/gatecrasher/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Context("when no relevant env vars are set", func() {
		BeforeEach(func() {
			Expect(os.Unsetenv("TARGET")).NotTo(HaveOccurred())
			Expect(os.Unsetenv("POLL_INTERVAL_IN_MS")).NotTo(HaveOccurred())
			Expect(os.Unsetenv("TOTAL_NUMBER_OF_REQUESTS")).NotTo(HaveOccurred())
			Expect(os.Unsetenv("REPORT_INTERVAL_IN_REQUESTS")).NotTo(HaveOccurred())
		})
		It("generates a config struct from defaults", func() {
			expectedConfig := config.Config{
				Target:                   "example.com",
				PollIntervalInMs:         1,
				TotalNumberOfRequests:    10,
				ReportIntervalInRequests: 5,
			}
			Expect(config.Load()).To(Equal(expectedConfig))
		})
	})
	Context("when all relevant env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("TARGET", "myurl")
			os.Setenv("POLL_INTERVAL_IN_MS", "10000")
			os.Setenv("TOTAL_NUMBER_OF_REQUESTS", "10000")
			os.Setenv("REPORT_INTERVAL_IN_REQUESTS", "500")
		})
		It("generates a config struct from env vars", func() {
			expectedConfig := config.Config{
				Target:                   "myurl",
				PollIntervalInMs:         10000,
				TotalNumberOfRequests:    10000,
				ReportIntervalInRequests: 500,
			}
			Expect(config.Load()).To(Equal(expectedConfig))
		})
	})
})
