package main

import (
	"flag"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type cfTeams struct {
	Teams []cfTeam `yaml:"teams"`
}

type cfTeam struct {
	Name                     string   `yaml:"name"`
	EnableSlackNotifications *bool    `yaml:"enable_slack_notifications"`
	Releases                 []string `yaml:"releases"`
	SlackChannel             string   `yaml:"slack_channel"`
}

func main() {
	var cfTeamsPath string
	flag.StringVar(&cfTeamsPath, "cf-teams", "", "Path to the CF Teams yaml file")

	var releaseRepository string
	flag.StringVar(&releaseRepository, "release-repository", "", "Release repository, as it appears on bosh.io (e.g. cloudfoundry/capi-release)")

	flag.Parse()

	if cfTeamsPath == "" {
		fmt.Fprintf(os.Stderr, "-cf-teams is a required flag\n")
		os.Exit(1)
	}

	if releaseRepository == "" {
		fmt.Fprintf(os.Stderr, "-release-repository is a required flag\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Parsing %s...\n", cfTeamsPath)
	teams, err := getTeams(cfTeamsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load teams: %s\n", err.Error())
		os.Exit(1)
	}

	var teamSlackChannel string

	fmt.Fprintf(os.Stderr, "Searching for team responsible for %s...\n", releaseRepository)
	for _, team := range teams.Teams {
		for _, release := range team.Releases {
			if release == releaseRepository {
				if team.EnableSlackNotifications != nil && !*team.EnableSlackNotifications {
					fmt.Fprintf(os.Stderr, "Found %s team. Slack notifications are disabled.\n", team.Name)
				} else {
					fmt.Fprintf(os.Stderr, "Found %s team. Slack notifications are enabled.\n", team.Name)
					teamSlackChannel = team.SlackChannel
				}

				break
			}
		}
	}

	fmt.Printf("%s", teamSlackChannel)
}

func getTeams(teamsPath string) (cfTeams, error) {
	var teams cfTeams
	yamlFile, err := os.ReadFile(teamsPath)
	if err != nil {
		return teams, err
	}

	err = yaml.Unmarshal(yamlFile, &teams)
	if err != nil {
		return teams, err
	}

	return teams, nil
}
