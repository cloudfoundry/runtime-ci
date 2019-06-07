package tracker_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		projectID = 0
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
			projectID = 12345

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
			expectedID, _ := trackerAPI.ListArgsForCall(0)
			Expect(expectedID).To(Equal(projectID), "expcted projectID value")
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
})
