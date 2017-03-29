# S3 Prober
A very simple Bash program
intended to measure transfer times to and from S3 
when pushed to Cloud Foundry.
The hope is that if there is
S3 service degradation 
that this will help reproduce that behavior
and make it visible in a Datadog dashboard.

## Pushing to Cloud Foundry
Push with the binary buildpack:

```
cf -f <path_to_manifest> push
```

This will push
with the specified `manifest.yml` file. Since this app
requires a S3 secret and Datadog API key, only a sample manifest is checked in.

Note that you will need to `cf login` first.

## Configuration
These environment variables are currently respected:
- `POLL_INTERVAL_IN_SECONDS` is the time between one batch of PUT, GET,
and DELETE.
- `TOTAL_NUMBER_OF_REQUESTS` is the number of requests
that will be made.
If it is set to 0
or less
it will make an unlimited number of requests.
- `S3_KEY` is the key to your AWS account.
- `S3_SECRET` is the secret to your AWS account.
- `S3_BUCKET` is the bucket to probe against.
- `DD_API_KEY` is your Datadog API key.
