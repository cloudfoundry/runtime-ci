package tracker

import (
	"gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

//go:generate counterfeiter . TrackerAPI
type TrackerAPI interface {
	List(projectID int, filter string) ([]*pivotal.Story, error)
}

type Client struct {
	api TrackerAPI

	projectID int
}

func NewClient(api TrackerAPI, projectID int) Client {
	return Client{api: api, projectID: projectID}
}

func (c Client) ScanForFlakeStory() (bool, error) {
	stories, err := c.api.List(c.projectID, "label:cats-flake-fix AND -state:accepted")
	if err != nil {
		return false, err
	}

	return len(stories) != 0, nil
}
