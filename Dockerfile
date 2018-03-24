#
# Commander
# Astronomer Platform Provisioning Service

FROM alpine:3.7
MAINTAINER Astronomer <humans@astronomer.io>

ENV REPO="github.com/astronomerio/commander"

ENV GCLOUD_VERSION="186.0.0"
ENV GCLOUD_FILE="google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz"
ENV GCLOUD_URL="https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/${GCLOUD_FILE}"

RUN apk update \
	&& apk add \
		build-base \
		ca-certificates \
		go \
		python2 \
	&& wget ${GCLOUD_URL} \
	&& tar -xvf ${GCLOUD_FILE} \
	&& rm ${GCLOUD_FILE} \
	&& google-cloud-sdk/install.sh \
	&& mkdir -p /opt \
	&& mv google-cloud-sdk /opt

WORKDIR /usr/lib/go/src/${REPO}

COPY . .

RUN make build
ENTRYPOINT ["./commander"]