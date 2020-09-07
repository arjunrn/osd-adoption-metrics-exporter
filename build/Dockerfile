FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/osd-adoption-metrics-exporter \
    USER_UID=1001 \
    USER_NAME=osd-adoption-metrics-exporter

# install operator binary
COPY build/_output/bin/osd-adoption-metrics-exporter ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
