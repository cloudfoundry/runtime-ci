# Gatecrasher
A very simple go program
intended to generate hairpin traffic
when pushed to Cloud Foundry.
The hope is that if there are
502 bad gateway errors
that happen more frequently
with requests that originate from within
the GCP project they're connecting to,
that this will help reproduce that behavior.

## Pushing to Cloud Foundry
The simplest way is to compile
and then push with the binary buildpack:

```
GOOS=linux GOARCH=amd64 go build main.go && cf push
```

This will push
with the `manifest.yml` file in this directory,
which specifies the appropriate buildpack
and start command.
It also arranges for the app to be routeless,
and to have no healthcheck,
as it's a "worker."

Note that you will need to `cf login` first.

## Configuration
These environment variables are currently respected:
- `POLL_INTERVAL_IN_MS` is the time between requests.
- `TARGET` is the full URL,
including protocol,
that requests are to be made against.
- `REPORT_INTERVAL_IN_REQUESTS` is the number of requests
that will be made between summary reports, which will
summarize the percentage of success requests.
- `TOTAL_NUMBER_OF_REQUESTS` is the number of requests
that will be made.
If it is set to 0
or less
it will make an unlimited number of requests.
`SKIP_INDIVIDUAL_REQUESTS` is the flag for suppressing
logs coming from each request, and it is default to false.
