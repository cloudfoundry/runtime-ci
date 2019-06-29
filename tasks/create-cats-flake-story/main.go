package main

import (
	"log"

	"github.com/cloudfoundry/runtime-ci/tasks/create-cats-flake-story/tracker"

	"github.com/spf13/viper"
	"gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

const (
	envProjectID = "project_id"
	envAPIToken  = "api_token"
)

func init() {
	viper.SetEnvPrefix("tracker")

	err := viper.BindEnv(
		// TRACKER_PROJECT_ID
		envProjectID,
		// TRACKER_API_TOKEN
		envAPIToken,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	apiToken := viper.GetString(envAPIToken)
	projectID := viper.GetInt(envProjectID)

	storyAPI := pivotal.NewClient(apiToken).Stories
	client := tracker.NewClient(storyAPI, projectID)

	exists, err := client.ScanForFlakeStory()
	if err != nil {
		log.Fatalln(err)
	}
	if exists {
		return
	}

	if err := client.CreateCATsFlakeStory(); err != nil {
		log.Fatalln(err)
	}
}
