package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
)

var (
	failLocation string
	terraformDir string
)

func init() {
	pflag.StringVarP(&failLocation, "fail-location", "l", "", "the failure location printed by ginkgo ")
	pflag.StringVarP(&terraformDir, "terraform-directory", "t", "", "the location where the terraform metadata is location")
}

func main() {
	pflag.Parse()

	if failLocation == "" {
		fmt.Println("fail-location required")
		pflag.Usage()
		fmt.Println(pflag.ErrHelp)
		os.Exit(1)
	}

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`select
	'https://release-integration.ci.cf-app.com/teams/main/pipelines/cf-deployment/jobs/stable-periodic-cats/builds/' || builds.name
	from build_events
	join builds on build_events.build_id=builds.id

	where
	type = 'log' and
	pipeline_id=9 and
	job_id=253
	and payload like '%' || $1 || '%'
	group by builds.name
	order by builds.name::integer desc`, failLocation)
	if err != nil {
		log.Fatal(err)
	}

	// var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)
	}
}

func connect() (*sql.DB, error) {
	initOut, err := exec.Command("terraform", "init", terraformDir).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing 'terraform init':\n%+v", string(initOut))
	}

	outputOut, err := exec.Command("terraform", "output", "--state", path.Join(terraformDir, "terraform.tfstate")).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing 'terraform output':\n%+v", string(outputOut))
	}

	host, password := parseTerraform(string(outputOut))
	pqOpts := fmt.Sprintf("host=%s user=atc dbname=concourse password=%s", host, password)
	db, err := sql.Open("postgres", pqOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func parseTerraform(output string) (host, password string) {
	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		parsed := strings.Split(line, " ")
		switch parsed[0] {
		case "concourse_db_password":
			password = parsed[2]
		case "db_host":
			host = parsed[2]
		}
	}

	return host, password
}
