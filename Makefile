OUTPUT ?= commander
IMAGE_NAME ?= astronomerinc/ap-${OUTPUT}
PROTO_SRC ?= $(shell pwd)/_proto
PROTO_DEST ?= $(shell pwd)/pkg/proto
.DEFAULT_GOAL := build

dep:
	dep ensure

build:
	CGO_ENABLED=0 go build -o ${OUTPUT} main.go

dependencies:
	dep ensure -vendor-only -v

build-image:
	docker build -t ${IMAGE_NAME} .

build-proto:
	mkdir -p ${PROTO_DEST}
	protoc --proto_path=${PROTO_SRC} --go_out=plugins=grpc:${PROTO_DEST} $(shell find ${PROTO_SRC} -type f -name "*.proto")

test:
	go test ./...

install: build
	mkdir -p $(DESTDIR)
	cp ${OUTPUT} $(DESTDIR)
