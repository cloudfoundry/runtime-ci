#!/bin/bash

# Utility that uses the cf cli's knowledge of request guids to search platform
# logs for useful information about what happened.  It's aimed at helping the
# Release Integration investigate the feasibility of implementing a low-effort
# trace for use in CATs.


# TODO: Clean up the output to make it more useful like this: `./what-hapened.sh | ag -v === | sort | uniq > wh.out`
# TODO: Maybe go!

# If we get killed, kill backgrounded processes
trap 'kill $(jobs -p) > /dev/null 2>&1' SIGTERM SIGINT

TMP_LOG=$(mktemp)
# TMP_LOG_CC=$(mktemp)
# TMP_LOG_ROUTER=$(mktemp)
echo "Main log is [ $TMP_LOG ]"
# echo "CC log is [ $TMP_LOG_CC ]"
# echo "ROUTER log is [ $TMP_LOG_ROUTER ]"

# TODO: Take cli command as argument to script
guids=$(cf -v s | ag X-Vcap-Request-Id | awk '{print $2}')

# I wrote a function so that we could keep execution time under control in
# case we want to fan out to grabbing several targeted logs.
grab_bosh_logs() {
  # bosh logs --job cloud_controller_ng --num 100 2> /dev/null > $TMP_LOG_CC &
  # bosh logs --job router --num 100 2> /dev/null > $TMP_LOG_ROUTER &
  bosh logs --num 1000 2> /dev/null > $TMP_LOG &

  for job in $(jobs -p); do
    wait "$job"
  done
}

grab_bosh_logs

# The guids that come back from the CLI often come back in the form
# <router-request-id>::<cf-cli-request-id>.  This splits them up to
# improve our chances with search.
good_guids=()
for guid in $guids; do
  if [[ $guid = *"::"* ]]; then
    g1=$(echo $guid | awk -F"::" '{print $1}')
    g2=$(echo $guid | awk -F"::" '{print $2}')

    good_guids+=("$g1")
    good_guids+=("$g2")
  else
    good_guids+=("$guid")
  fi
done

# `good_guids` is a proper bash array, which requires this weird
# access syntax.
for guid in "${good_guids[@]}"; do
  echo "=== looking for guid [ $guid ] in Main log ==="
  ag $guid $TMP_LOG
  echo

  # echo "=== looking for guid [ $guid ] in CC log ==="
  # ag $guid $TMP_LOG_CC
  # echo
  #
  # echo "=== looking for guid [ $guid ] in Router log ==="
  # ag $guid $TMP_LOG_ROUTER
  # echo
done
