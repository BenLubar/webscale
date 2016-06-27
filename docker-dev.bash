#!/bin/bash -ex

# The first 3 octets of an IPv4 /24 dedicated to #webscale. Defaults to 172.19.1
WEBSCALE_SUBNET=${WEBSCALE_SUBNET:-172.19.1}

./test-docker.bash

if [[ "$1" = "--reset" ]]; then
	docker pull benlubar/webscale:postgres
	docker stop webscale-postgres webscale-app-{0..1}
	docker rm -v webscale-postgres
	docker run -d --name webscale-postgres --restart unless-stopped --net webscale --ip "$WEBSCALE_SUBNET.2" benlubar/webscale:postgres
fi

docker build -t benlubar/webscale .

for i in {0..1}; do
	docker stop webscale-app-$i
	docker rm -v webscale-app-$i
	docker run -d --name webscale-app-$i --restart unless-stopped --net webscale --ip "$WEBSCALE_SUBNET.25$i" benlubar/webscale
done
