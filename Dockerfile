#
# Commander
# Astronomer Platform Provisioning Service

FROM golang:1.9.2-alpine3.7
MAINTAINER Astronomer <humans@astronomer.io>

ENV GCLOUD_VERSION="185.0.0"
ENV GCLOUD_FILE="google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz"
ENV GCLOUD_URL="https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/${GCLOUD_FILE}"

WORKDIR /opt
RUN set -x \
	&& apk update \
	&& apk add --no-cache --virtual .build-deps \
		build-base \
	&& apk --no-cache add \
		ca-certificates \
		python2 \
	&& wget ${GCLOUD_URL} \
	&& tar -xvf ${GCLOUD_FILE} \
	&& rm ${GCLOUD_FILE} \
	&& google-cloud-sdk/install.sh

WORKDIR /go/src/github.com/astronomerio/commander
COPY . .
RUN make build

ENTRYPOINT ["./commander"]
