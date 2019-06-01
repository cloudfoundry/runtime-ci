package tracker

import (
	"strings"

	"gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

//go:generate counterfeiter . TrackerAPI
type TrackerAPI interface {
	List(projectID int, filter string) ([]*pivotal.Story, error)
}

type Client struct {
	api TrackerAPI

	projectID string
}

func NewClient(api TrackerAPI, projectID string) Client {
	return Client{api: api, projectID: projectID}
}

func (c Client) ScanForFlakeStory() bool {
	stories, _ := c.api.List(88, "")

	for _, story := range stories {
		if story.State == "accepted" {
			continue
		}

		if strings.HasPrefix(story.Name, "CAT Failure Fix") {
			return true
		}
	}
	return false
}
