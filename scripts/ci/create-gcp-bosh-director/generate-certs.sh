#!/bin/bash
# This script adapted from the one presented at https://bosh.io/docs/director-certs.html on June 11 2016

set -e

name=$1
ip=$2
output_dir=$3

# rm -rf $certs && mkdir -p $certs
pushd $output_dir

	if [ ! -e certs ]; then
		mkdir certs
	fi
	cd certs

	echo "Generating new CA..."
	openssl genrsa -out rootCA.key 2048
	yes "" | openssl req -x509 -new -nodes -key rootCA.key \
		-out rootCA.pem -days 99999


	cat >openssl-exts.conf <<-EOL
	extensions = san
	[san]
	subjectAltName = IP:${ip}
	EOL


	echo "Generating certificate signing request for ${ip}..."
	# golang requires to have SAN for the IP
	openssl req -new -nodes -newkey rsa:2048 \
		-out ${name}.csr -keyout ${name}.key \
		-subj "/C=US/O=BOSH/CN=${ip}"


	echo "Generating certificate ${ip}..."
	openssl x509 -req -in ${name}.csr \
		-CA rootCA.pem -CAkey rootCA.key -CAcreateserial \
		-out ${name}.crt -days 99999 \
		-extfile ./openssl-exts.conf

	echo "Deleting certificate signing request and config..."
	rm ${name}.csr
	rm ./openssl-exts.conf

	echo "Finished..."
	ls -la .

popd
