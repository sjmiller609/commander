OUTPUT ?= commander
IMAGE_NAME ?= astronomerinc/ap-pce-${OUTPUT}

.DEFAULT_GOAL := build

dep:
	dep ensure

build: dep
	go build -o ${OUTPUT} main.go

build-image:
	docker build -t ${IMAGE_NAME} .

install: build
	mkdir -p $(DESTDIR)
	cp ${OUTPUT} $(DESTDIR)
