FROM alpine

ARG DOCKER_BINARY

ADD $DOCKER_BINARY /opt/bin/konfigurator

ENTRYPOINT ["/opt/bin/konfigurator"]