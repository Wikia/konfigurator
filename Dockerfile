FROM scratch

ARG DOCKER_BINARY

ADD $DOCKER_BINARY /opt/bin/konfigurator

ENTRYPOINT ["/opt/bin/konfigurator"]

VOLUME /config
ENV HOME="/config"
WORKDIR /config

RUN chown 65534:65534 /config

USER 65534:65534
