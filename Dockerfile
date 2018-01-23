#
# Commander
# Astronomer Platform Provisioning Service

FROM golang:1.9.2
MAINTAINER Astronomer <humans@astronomer.io>

WORKDIR /go/src/github.com/astronomerio/commander
COPY . .
RUN make build

EXPOSE 8881
ENTRYPOINT ["./commander"]
