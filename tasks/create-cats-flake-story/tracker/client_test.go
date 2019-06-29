package tracker_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"

	. "github.com/cloudfoundry/runtime-ci/tasks/create-cats-flake-story/tracker"
	fakes "github.com/cloudfoundry/runtime-ci/tasks/create-cats-flake-story/tracker/trackerfakes"
)

var _ = Describe("Client", func() {
	var (
		trackerAPI *fakes.FakeTrackerAPI
		projectID  int

		client Client
	)

	BeforeEach(func() {
		trackerAPI = new(fakes.FakeTrackerAPI)
		projectID = 12345
	})

	JustBeforeEach(func() {
		client = NewClient(trackerAPI, projectID)
	})

	Describe("ScanForFlakeStory", func() {
		var (
			returnStories []*pivotal.Story
			returnErr     error

			actualExists bool
			actualErr    error
		)

		BeforeEach(func() {
			returnStories = nil

			actualExists = false
			actualErr = nil
		})

		JustBeforeEach(func() {
			trackerAPI.ListReturns(returnStories, returnErr)

			actualExists, actualErr = client.ScanForFlakeStory()
		})

		It("uses the projectID query the API", func() {
			Expect(trackerAPI.ListCallCount()).To(Equal(1), "expected call count")
			actualID, _ := trackerAPI.ListArgsForCall(0)
			Expect(actualID).To(Equal(projectID), "expcted projectID value")
		})

		It("filters by 'label:cats-flake-fix AND -state:accepted'", func() {
			Expect(trackerAPI.ListCallCount()).To(Equal(1), "expected call count")
			_, filter := trackerAPI.ListArgsForCall(0)
			Expect(filter).To(Equal("label:cats-flake-fix AND -state:accepted"), "expected filter value")
		})

		Context("when no stories are returned", func() {
			BeforeEach(func() {
				returnStories = nil
			})

			It("returns false for existing CATs flake", func() {
				Expect(actualExists).To(BeFalse())
			})
		})

		Context("when at least one story is returned", func() {
			BeforeEach(func() {
				returnStories = append(returnStories, &pivotal.Story{Name: "Some Tracker Story"})
			})

			It("returns true for existing CATs flake", func() {
				Expect(actualExists).To(BeTrue(), "expected story to exist")
			})
		})

		Context("when the api returns an error", func() {
			BeforeEach(func() {
				returnErr = errors.New("a crazy error")
			})

			It("returns an error", func() {
				Expect(actualErr).To(MatchError(returnErr), "expected an error")
			})
		})
	})

	Describe("CreateCATsFlakeStory", func() {
		var (
			returnStory *pivotal.Story
			returnResp  *http.Response
			returnErr   error

			actualErr error
		)

		BeforeEach(func() {
			returnStory = nil
			returnResp = &http.Response{
				Body: ioutil.NopCloser(strings.NewReader("Fake body")),
			}
			returnErr = nil
		})

		JustBeforeEach(func() {
			trackerAPI.CreateReturns(returnStory, returnResp, returnErr)

			actualErr = client.CreateCATsFlakeStory()
		})

		It("uses the projectID query the API", func() {
			Expect(trackerAPI.CreateCallCount()).To(Equal(1), "expected call count")
			actualID, _ := trackerAPI.CreateArgsForCall(0)
			Expect(actualID).To(Equal(projectID), "expcted projectID value")
		})

		It("creates the expected CATs Flake Story", func() {
			Expect(trackerAPI.CreateCallCount()).To(Equal(1), "expected call count")
			_, actualStory := trackerAPI.CreateArgsForCall(0)
			Expect(*actualStory).To(MatchFields(IgnoreExtras, Fields{
				"Name": Equal("CAT Failure Fix -- [Unstarted]"),
				"Description": Equal(`**Process:**
1. Determine if the the incidence rate of any flake is high enough to try to fix
  - [base honeycomb query](https://ui.honeycomb.io/cf-release-integration/datasets/canonical-cats?query={"breakdowns":["Description","State"],"calculations":[{"op":"COUNT"}],"filters":[{"column":"State","op":"=","value":"failed"}],"orders":[{"op":"COUNT","order":"descending"}],"time_range":1209600})
1. Create honeycomb query for the particular failure and the incidence rate at the time of starting the story.
  - Change the bracketed portion of this story title to reflect the failure description.
	- [honeycomb url] (PASTE URL HERE)
	- INCIDENCE RATE/14 days as of DATE HERE
1. Analyze the failure and come up with an approaceh to reduce/eliminate the error rate.
  - Document the approach as a comment in this story
`),
				"Type":  Equal("chore"),
				"State": Equal("planned"),
			}))
		})

		Context("when the request succeeds", func() {
			It("successfully returns", func() {})
		})

		Context("when the request fails", func() {
			Context("do to a a failure to create the request", func() {
				BeforeEach(func() {
					returnErr = errors.New("oh no banana")
				})

				It("returns the error", func() {
					Expect(actualErr).To(MatchError("invalid request: oh no banana"))
				})
			})

			Context("do to a bad response", func() {
				BeforeEach(func() {
					returnResp.StatusCode = http.StatusTeapot
					returnResp.Status = "418 I'm a teapot"
				})

				It("returns the error", func() {
					Expect(actualErr).To(MatchError("invalid response: 418 I'm a teapot"))
				})
			})
		})
	})
})
