#
# Commander
# Astronomer Platform Provisioning Service

FROM golang:latest AS build-env

MAINTAINER Astronomer <humans@astronomer.io>

ARG BUILD_NUMBER=-1
LABEL io.astronomer.docker.build.number=$BUILD_NUMBER
LABEL io.astronomer.docker.module="astronomer"
LABEL io.astronomer.docker.component="commander"
LABEL io.astronomer.docker.environment="development"

ENV REPO="github.com/astronomerio/commander"
ENV GOPATH=/root/go
ENV GOBIN=/root/go/bin
ENV PATH=${PATH}:${GOBIN}
ENV PROTOC_VERSION="3.10.0"
ENV PROTOC_FILE="protoc-${PROTOC_VERSION}-linux-x86_64"
ENV PROTOC_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_FILE}.zip"
ENV PROTOC_GEN_GO_VERSION="v1.0.0"

RUN apt-get update && \
    apt-get install -y dtrx

RUN wget $PROTOC_URL \
    && dtrx ./${PROTOC_FILE}.zip \
    && mv ./$PROTOC_FILE/bin/protoc /usr/local/sbin/protoc

RUN go get -u github.com/golang/dep/cmd/dep

# Switch into source dir under GOPATH
WORKDIR ${GOPATH}/src/${REPO}

COPY Makefile Gopkg.toml Gopkg.lock ./

RUN make dependencies
RUN go get -d -u github.com/golang/protobuf/protoc-gen-go \
    && git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $PROTOC_GEN_GO_VERSION \
    && go install github.com/golang/protobuf/protoc-gen-go

COPY . ./

# Build commander
RUN make build-proto
RUN make build && cp ./commander /commander

FROM alpine:latest

WORKDIR /app
COPY --from=build-env /commander /app/

# Run Commander
ENTRYPOINT /app/commander
