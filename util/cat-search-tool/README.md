When provided a ginkgo error location and a valid concourse terraform directory, 
returns a list of urls that link back to Bellatrix failures in our Concourse pipeline.
```
$: cd ${runtime-ci}/util/cat-search-tool
$: go run main.go -l <failure-location line from honeycomb> -t <concourse-env/terraform/cloudsql>
```

May want to consider extending this simple implementation so that it:
1. can be used to search issues associated with other pipelines (right now it's hardcoded to Bellatrix)
