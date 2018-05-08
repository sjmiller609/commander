#
# Commander
# Astronomer Platform Provisioning Service

FROM astronomerinc/ap-base
MAINTAINER Astronomer <humans@astronomer.io>

ARG BUILD_NUMBER=-1
LABEL io.astronomer.docker.build.number=$BUILD_NUMBER
LABEL io.astronomer.docker.module="core"
LABEL io.astronomer.docker.component="houston"

# Update apk
RUN apk update

ENV REPO="github.com/astronomerio/commander"

ENV GCLOUD_VERSION="186.0.0"
ENV GCLOUD_FILE="google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz"
ENV GCLOUD_URL="https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/${GCLOUD_FILE}"
ENV GOPATH=/root/go
ENV GOBIN=/root/go/bin
ENV PATH=${PATH}:${GOBIN}

RUN apk add \
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

WORKDIR ${GOPATH}/src/${REPO}

COPY . .

RUN make build
ENTRYPOINT ["./commander"]