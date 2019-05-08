package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type cfTeams struct {
	Teams []cfTeam `yaml:"teams"`
}

type cfTeam struct {
	Name         string   `yaml:"name"`
	SlackChannel string   `yaml:"slack_channel"`
	Releases     []string `yaml:"releases"`
}

func main() {
	var cfTeamsPath string
	flag.StringVar(&cfTeamsPath, "cf-teams", "", "path to the CF Teams yaml file")
	var repository string
	flag.StringVar(&repository, "repository", "", "repository, as it appears on bosh.io (e.g. cloudfoundry/capi-release)")
	flag.Parse()

	if cfTeamsPath == "" {
		fmt.Fprintf(os.Stderr, "-cf-teams is a required flag\n")
		os.Exit(1)
	}
	if repository == "" {
		fmt.Fprintf(os.Stderr, "-repository is a required flag\n")
		os.Exit(1)
	}

	teams, err := getTeams(cfTeamsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load teams: %s\n", err.Error())
		os.Exit(1)
	}

	var teamSlackChannel string

	for _, team := range teams.Teams {
		for _, release := range team.Releases {
			if release == repository {
				teamSlackChannel = team.SlackChannel
			}
		}
	}

	fmt.Printf("%s", teamSlackChannel)
}

func getTeams(teamsPath string) (cfTeams, error) {
	var teams cfTeams
	yamlFile, err := ioutil.ReadFile(teamsPath)
	if err != nil {
		return teams, err
	}

	err = yaml.Unmarshal(yamlFile, &teams)
	if err != nil {
		return teams, err
	}

	return teams, nil
}
