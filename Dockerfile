#
# Copyright 2018 Astronomer Inc.
#
# Licensed under the Apache License, Version 3.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM astronomerinc/ap-base:0.10.3 AS build-env

ENV REPO="github.com/astronomerio/commander"
ENV GOPATH=/root/go
ENV GOBIN=${GOPATH}/bin
ENV PATH=${PATH}:${GOBIN}

WORKDIR /usr/lib/go/src/${REPO}

RUN apk update \
	  && apk add --no-cache --virtual .build-deps \
	  	build-base \
	  	go

COPY . .

RUN make DESTDIR=/usr/bin install

RUN mkdir -p "${GOPATH}/src" \
	  && mkdir -p "${GOPATH}/bin" \
	  && rm -rf /usr/lib/go \
	  && apk del .build-deps

ENTRYPOINT ["commander"]

#
# Build final image
#
FROM astronomerinc/ap-base:0.10.3
MAINTAINER Astronomer <humans@astronomer.io>

ARG BUILD_NUMBER=-1
LABEL maintainer="Astronomer <humans@astronomer.io>"
LABEL io.astronomer.docker.build.number=$BUILD_NUMBER
LABEL io.astronomer.docker.module="astronomer"
LABEL io.astronomer.docker.component="commander"
LABEL io.astronomer.docker.environment="development"

WORKDIR /app
COPY --from=build-env /usr/bin/commander /app/

# Run Commander
ENTRYPOINT /app/commander
