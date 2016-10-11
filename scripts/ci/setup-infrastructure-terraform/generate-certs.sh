#!/bin/bash
# This script adapted from the one presented at https://bosh.io/docs/director-certs.html on June 11 2016

set -e

certs=`dirname $0`/certs
name=$1
ip=$2
ip_filename=$3
output_dir=$4

# rm -rf $certs && mkdir -p $certs
pushd $4

	if [ ! -e certs ]; then
		mkdir certs
	fi
	cd $certs

	echo "Generating new CA..."
	openssl genrsa -out rootCA-${ip_filename}.key 2048
	yes "" | openssl req -x509 -new -nodes -key rootCA-${ip_filename}.key \
		-out rootCA-${ip_filename}.pem -days 99999


	cat >openssl-exts.conf <<-EOL
	extensions = san
	[san]
	subjectAltName = IP:${ip}
	EOL


	echo "Generating certificate signing request for ${ip}..."
	# golang requires to have SAN for the IP
	openssl req -new -nodes -newkey rsa:2048 \
		-out ${name}-${ip_filename}.csr -keyout ${name}-${ip_filename}.key \
		-subj "/C=US/O=BOSH/CN=${ip}"


	echo "Generating certificate ${ip}..."
	openssl x509 -req -in ${name}-${ip_filename}.csr \
		-CA rootCA-${ip_filename}.pem -CAkey rootCA-${ip_filename}.key -CAcreateserial \
		-out ${name}-${ip_filename}.crt -days 99999 \
		-extfile ./openssl-exts.conf

	echo "Deleting certificate signing request and config..."
	rm ${name}-${ip_filename}.csr
	rm ./openssl-exts.conf

	echo "Finished..."
	ls -la .

popd
