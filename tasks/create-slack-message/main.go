package main

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type cfTeams struct {
	teams []cfTeam `yaml:"teams"`
}

type cfTeam struct {
	name         string   `yaml:"name"`
	slackChannel string   `yaml:"slackChannel"`
	releases     []string `yaml:"releases"`
}

func main() {
	releaseName := os.Args[2]

	teams := getTeams(os.Args[1])

}

func getTeams(teamsPath string) cfTeams {
	teams := cfTeams{}
	yamlFile, err := ioutil.ReadFile(teamsPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &teams)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return teams
}
