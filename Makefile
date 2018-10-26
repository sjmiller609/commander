OUTPUT ?= commander
IMAGE_NAME ?= astronomerinc/ap-${OUTPUT}
PROTO_SRC ?= $(shell pwd)/_proto
PROTO_DEST ?= $(shell pwd)/pkg/proto
.DEFAULT_GOAL := build

dep:
	dep ensure

build:
	go build -o ${OUTPUT} main.go

build-image:
	docker build -t ${IMAGE_NAME} .

build-proto:
	protoc --proto_path=${PROTO_SRC} --go_out=plugins=grpc:${PROTO_DEST} $(shell find ${PROTO_SRC} -type f -name "*.proto")

test:
	go test ./...

install: build
	mkdir -p $(DESTDIR)
	cp ${OUTPUT} $(DESTDIR)
