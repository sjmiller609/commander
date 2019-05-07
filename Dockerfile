#
# Commander
# Astronomer Platform Provisioning Service

FROM alpine:3.9
MAINTAINER Astronomer <humans@astronomer.io>

ARG BUILD_NUMBER=-1
LABEL io.astronomer.docker.build.number=$BUILD_NUMBER
LABEL io.astronomer.docker.module="astronomer"
LABEL io.astronomer.docker.component="houston"
LABEL io.astronomer.docker.environment="development"

ENV REPO="github.com/astronomerio/commander"
ENV GCLOUD_VERSION="186.0.0"
ENV GCLOUD_FILE="google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz"
ENV GCLOUD_URL="https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/${GCLOUD_FILE}"
ENV GOPATH=/root/go
ENV GOBIN=/root/go/bin
ENV PATH=${PATH}:${GOBIN}

# Install dependencies
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
	&& mv google-cloud-sdk /opt \
	&& mkdir -p "${GOPATH}" \
	&& mkdir "${GOPATH}/src" \
	&& mkdir "${GOPATH}/bin"

# Switch into source dir under GOPATH
WORKDIR ${GOPATH}/src/${REPO}

# Copy all source code in
COPY . .

# Run Commander
ENTRYPOINT ["go", "run", "main.go"]
