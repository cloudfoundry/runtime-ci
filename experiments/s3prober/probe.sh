#!/bin/bash -eu

function send_datadog() {
  echo "Posting $1 result to DataDog ..."
  currenttime=$(date +%s)
  curl  -X POST -H "Content-type: application/json" \
  -d "{ \"series\" :
           [{\"metric\":\"test.metric.s3prober.$1\",
            \"points\":[[$currenttime, $2]],
            \"type\":\"gauge\",
            \"host\":\"ord.example.com\",
            \"tags\":[\"environment:test\"]}
          ]
      }" \
  "https://app.datadoghq.com/api/v1/series?api_key=$DD_API_KEY"
  echo
  echo
}

function putS3() {
  echo "Running put test to S3 ..."
  path=$1
  file=$2
  aws_path="/"
  date=$(date +"%a, %d %b %Y %T %z")
  acl="x-amz-acl:public-read"
  content_type='application/x-compressed-tar'
  string="PUT\n\n${content_type}\n${date}\n${acl}\n/${S3_BUCKET}${aws_path}${file}"
  signature=$(echo -en "${string}" | openssl sha1 -hmac "${S3_SECRET}" -binary | base64)
  curl -X PUT -T "${path}/${file}" \
    -H "Host: ${S3_BUCKET}.s3.amazonaws.com" \
    -H "Date: $date" \
    -H "Content-Type: $content_type" \
    -H "$acl" \
    -H "Authorization: AWS ${S3_KEY}:${signature}" \
    "https://${S3_BUCKET}.s3.amazonaws.com${aws_path}${file}"
  echo $?
}

function deleteS3() {
  echo "Running delete test to S3 ..."
  file=$1
  aws_path="/"
  date=$(date +"%a, %d %b %Y %T %z")
  string="DELETE\n\n\n${date}\n/${S3_BUCKET}${aws_path}${file}"
  signature=$(echo -en "${string}" | openssl sha1 -hmac "${S3_SECRET}" -binary | base64)
  curl -X DELETE \
    -H "Host: ${S3_BUCKET}.s3.amazonaws.com" \
    -H "Date: $date" \
    -H "Authorization: AWS ${S3_KEY}:${signature}" \
    "https://${S3_BUCKET}.s3.amazonaws.com${aws_path}${file}"
}

function getS3() {
  echo "Running get test from S3 ..."
  file=$1
  aws_path="/"
  curl "https://${S3_BUCKET}.s3.amazonaws.com${aws_path}${file}" > "/tmp/${file}"
}

function main() {
  for (( i=0; ; i++ ))
  do
    if (( i > TOTAL_NUMBER_OF_REQUESTS )) ; then
      if (( TOTAL_NUMBER_OF_REQUESTS != 0 )) && (( TOTAL_NUMBER_OF_REQUESTS != -1 )); then
        exit 0
      fi
    fi
    starttime=$(($(date +%s%N)/1000000))

    # PUT to S3
    putS3 "./" "big_dora.tgz"

    endtime=$(($(date +%s%N)/1000000))

    # Send PUT time to datadog
    send_datadog put $(( endtime - starttime ))

    starttime=$(($(date +%s%N)/1000000))

    # GET from S3
    getS3 "big_dora.tgz"

    endtime=$(($(date +%s%N)/1000000))

    # Send GET time to datadog
    send_datadog get $(( endtime - starttime ))

    starttime=$(($(date +%s%N)/1000000))

    # DELETE from S3
    deleteS3 "big_dora.tgz"

    endtime=$(($(date +%s%N)/1000000))

    # Send GET time to datadog
    send_datadog delete $(( endtime - starttime ))

    sleep $(( POLL_INTERVAL_IN_SECONDS ))
  done
}

main

