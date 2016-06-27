FROM postgres:9.4

RUN apt-get update \
&& apt-get install -y postgresql-plperl-$PG_MAJOR=$PG_VERSION \
&& rm -rf /var/lib/apt/lists/*
