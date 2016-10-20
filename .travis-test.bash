#!/bin/bash

case "$DATABASE" in
postgres)
	data_source_name="dbname=travis_ci_test user=postgres sslmode=disable"
	go_tags=""
	;;
esac

rm -f all.prof
touch all.prof

go build -tags "$go_tags" ./cmd/webscale-upgrade && ./webscale-upgrade -db "$data_source_name" || exit $?

error_status=0

go list ./... | while read -r pkg; do
	go test -tags "$go_tags" -race -coverprofile this.prof -coverpkg ./... -v "$pkg" -db "$data_source_name" |& grep -v 'warning: no packages being tested depend on '
	error_status=$(( $error_status + $? ))

	test -f this.prof && gocovmerge this.prof all.prof > merged.prof && mv merged.prof all.prof
	rm -f this.prof
done

goveralls -coverprofile all.prof

rm -f all.prof merged.prof

exit $error_status
