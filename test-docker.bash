#!/bin/bash

# test-docker.bash tests runs webscale tests using docker to quickly create and
# clean up supported databases. It generates a coverage report for each run.

# The first 3 octets of an IPv4 /24 dedicated to #webscale. Defaults to 172.19.1
WEBSCALE_SUBNET=${WEBSCALE_SUBNET:-172.19.1}

db="host=$WEBSCALE_SUBNET.42 user=postgres sslmode=disable"

pg_versions="9.4"

for pg_version in $pg_versions; do
	# update the postgres image from the docker hub
	docker pull benlubar/webscale:postgres-$pg_version
done

for pg_version in $pg_versions; do
	# start the postgres server
	docker run -d --name webscale-postgres-test --net webscale --ip "$WEBSCALE_SUBNET.42" benlubar/webscale:postgres-$pg_version

	# wait for postgres to be ready
	until nc -z "$WEBSCALE_SUBNET.42" 5432; do sleep 1; done

	rm -f all.prof
	touch all.prof

	# run the upgrade script
	go build ./cmd/webscale-upgrade && ./webscale-upgrade -db "$db"
	error_status=$?
	rm -f webscale-upgrade

	# for each package, run the tests and add the coverage to all.prof
	go list ./... | while read -r pkg; do
		go test -race -coverprofile this.prof -coverpkg ./... -v "$pkg" -db "$db" |& grep -v 'warning: no packages being tested depend on '
		error_status=$(( $error_status + $? ))

		test -f this.prof && gocovmerge this.prof all.prof > merged.prof && mv merged.prof all.prof
		rm -f this.prof
	done

	# delete the postgres server
	docker stop webscale-postgres-test && docker rm -v webscale-postgres-test

	# write the coverage output
	go tool cover -html all.prof -o "coverage-pg$pg_version.html"
	rm -f all.prof merged.prof

	# if we failed, exit
	if [[ "$exit_status" -ne 0 ]]; then
		echo "One or more tests failed. Stopping."
		exit "$exit_status"
	fi
done
