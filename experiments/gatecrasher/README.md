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
