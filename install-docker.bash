#!/bin/bash -ex

# The first 3 octets of an IPv4 /24 dedicated to #webscale. Defaults to 172.19.1
WEBSCALE_SUBNET=${WEBSCALE_SUBNET:-172.19.1}

# update the postgres image from the docker hub. this is the postgres:latest
# image from Docker's official repository with the plperl extension added.
docker pull benlubar/webscale:postgres

# download or update the webscale image as well.
docker pull benlubar/webscale

# create a virtual network.
docker network create --subnet="$WEBSCALE_SUBNET.0/24" webscale

# start the postgres server.
docker run -d --name webscale-postgres --restart unless-stopped --net webscale --ip "$WEBSCALE_SUBNET.2" benlubar/webscale:postgres

# start two webscale servers so we can do a graceful restart.
for i in {0..1}; do
	docker run -d --name webscale-app-$i --restart unless-stopped --net webscale --ip "$WEBSCALE_SUBNET.25$i" benlubar/webscale
done
