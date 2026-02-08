FROM alpine:3
ARG TARGETPLATFORM

RUN apk add --no-cache \
  bash \
  traceroute \
  iputils-ping

COPY $TARGETPLATFORM/dping_exporter /usr/bin/

CMD ["/usr/bin/dping_exporter"]
