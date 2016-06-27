#!/bin/bash

# The first 3 octets of an IPv4 /24 dedicated to #webscale. Defaults to 172.19.1
WEBSCALE_SUBNET=${WEBSCALE_SUBNET:-172.19.1}

pg_versions=9.4 9.5 9.6

for pg_version in $pg_versions; do
	# update the postgres image from the docker hub
	docker pull benlubar/webscale:postgres-$pg_version
done

for pg_version in $pg_versions; do
	# start the postgres server
	docker run -d --name webscale-postgres-test --net webscale --ip "$WEBSCALE_SUBNET.42" benlubar/webscale:postgres-$pg_version

	# wait for postgres to be ready
	until nc -z "$WEBSCALE_SUBNET.42" 5432; do sleep 1; done

	# run the schema test first so if the database init fails it fails here
	go test -cover ./db/internal/schema -db "host=$WEBSCALE_SUBNET.42 user=postgres sslmode=disable" && \
	go test -cover ./... -db "host=$WEBSCALE_SUBNET.42 user=postgres sslmode=disable"

	exit_status=$?

	# delete the postgres server
	docker stop webscale-postgres-test && docker rm -v webscale-postgres-test

	if [[ "$exit_status" -ne 0 ]]; then
		exit "$exit_status"
	fi
done
