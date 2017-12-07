FROM alpine:3.6

ARG DOCKER_BINARY

ADD $DOCKER_BINARY /usr/bin/konfigurator

ENV HOME=/config
RUN adduser konfigurator -D -h /config
RUN chown -R konfigurator /config

VOLUME ["/config"]
USER konfigurator
