IMAGE_NAME ?= astronomerinc/ap-commander
OUTPUT ?= commander

build:
	go build -o ${OUTPUT} main.go

build-image:
	docker build -t ${IMAGE_NAME} .

install: build
	mkdir -p $(DESTDIR)
	cp ${OUTPUT} $(DESTDIR)
