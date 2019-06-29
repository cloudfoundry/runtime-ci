package tracker

import (
	"fmt"
	"net/http"

	"gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

//go:generate counterfeiter . API
type API interface {
	Create(projectID int, story *pivotal.StoryRequest) (*pivotal.Story, *http.Response, error)
	List(projectID int, filter string) ([]*pivotal.Story, error)
}

type Client struct {
	api API

	projectID int
}

func NewClient(api API, projectID int) Client {
	return Client{api: api, projectID: projectID}
}

func (c Client) ScanForFlakeStory() (bool, error) {
	stories, err := c.api.List(c.projectID, "label:cats-flake-fix AND -state:accepted")
	if err != nil {
		return false, err
	}

	return len(stories) != 0, nil
}

func (c Client) CreateCATsFlakeStory() error {
	_, resp, err := c.api.Create(c.projectID, newStoryRequest())
	if err != nil {
		return fmt.Errorf("invalid request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		return fmt.Errorf("invalid response: %s", resp.Status)
	}

	return nil
}

func newStoryRequest() *pivotal.StoryRequest {
	return &pivotal.StoryRequest{
		Name: "CAT Failure Fix -- [Unstarted]",
		Description: `**Process:**
1. Determine if the the incidence rate of any flake is high enough to try to fix
  - [base honeycomb query](https://ui.honeycomb.io/cf-release-integration/datasets/canonical-cats?query={"breakdowns":["Description","State"],"calculations":[{"op":"COUNT"}],"filters":[{"column":"State","op":"=","value":"failed"}],"orders":[{"op":"COUNT","order":"descending"}],"time_range":1209600})
1. Create honeycomb query for the particular failure and the incidence rate at the time of starting the story.
  - Change the bracketed portion of this story title to reflect the failure description.
	- [honeycomb url] (PASTE URL HERE)
	- INCIDENCE RATE/14 days as of DATE HERE
1. Analyze the failure and come up with an approaceh to reduce/eliminate the error rate.
  - Document the approach as a comment in this story
`,
		Type:  "chore",
		State: "planned",
	}
}
