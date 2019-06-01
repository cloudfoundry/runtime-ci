package tracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"

	. "github.com/cloudfoundry/runtime-ci/tasks/create-cats-flake-story/tracker"
	fakes "github.com/cloudfoundry/runtime-ci/tasks/create-cats-flake-story/tracker/trackerfakes"
)

var _ = Describe("Client", func() {
	var (
		trackerAPI *fakes.FakeTrackerAPI
		projectID  string

		client Client
	)

	BeforeEach(func() {
		trackerAPI = new(fakes.FakeTrackerAPI)
	})

	JustBeforeEach(func() {
		client = NewClient(trackerAPI, projectID)
	})

	Describe("ScanForFlakeStory", func() {
		var (
			stories []*pivotal.Story

			actualExists bool
		)

		BeforeEach(func() {
			stories = nil
			actualExists = false
		})

		JustBeforeEach(func() {
			trackerAPI.ListReturns(stories, nil)

			actualExists = client.ScanForFlakeStory()
		})

		It("uses the projectID query the API", func() {
			Expect(trackerAPI.ListCallCount()).To(Equal(1), "expected call count")
		})

		Context("when no stories are returned", func() {
			BeforeEach(func() { stories = nil })
			It("returns false for existing CATs flake", func() { Expect(actualExists).To(Equal(false)) })
		})

		Context("when at least one story is returned", func() {
			BeforeEach(func() {
				stories = append(stories, &pivotal.Story{Name: "DOG Failure Fix"})
			})

			Context("when a flake story does exist", func() {
				var story *pivotal.Story

				BeforeEach(func() {
					story = &pivotal.Story{Name: "CAT Failure Fix"}
					stories = append(stories, story)
				})

				Context("when the flake is in the accepted state", func() {
					BeforeEach(func() { story.State = "accepted" })
					It("returns false for existing CATs flake", func() { Expect(actualExists).To(Equal(false)) })
				})

				Context("when the flake is in any other state", func() {
					BeforeEach(func() { story.State = "some state" })
					It("returns true for existing CATs flake", func() { Expect(actualExists).To(Equal(true)) })
				})
			})

			Context("when a flake story doesn't exist", func() {
				It("returns false for existing CATs flake", func() { Expect(actualExists).To(Equal(false)) })
			})
		})
	})
})
