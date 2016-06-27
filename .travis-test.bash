#!/bin/bash

# Test the schema package first so the database init fails the correct test
go test -race -coverprofile all.prof -coverpkg ./... -v ./db/internal/schema -db 'dbname=travis_ci_test user=postgres sslmode=disable' |& grep -v 'warning: no packages being tested depend on '
error_status=$?

# Test all the packages and upload the coverage information
go list ./... | while read -r pkg; do
	go test -race -coverprofile this.prof -coverpkg ./... -v "$pkg" -db 'dbname=travis_ci_test user=postgres sslmode=disable' |& grep -v 'warning: no packages being tested depend on '
	error_status=$(( $error_status + $? ))

	cat this.prof >> all.prof 2> /dev/null
done

sed '1!{/^mode:/d;}' -i all.prof
goveralls -coverprofile all.prof

exit $error_status
