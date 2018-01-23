IMAGE_NAME ?= astronomerinc/ap-commander

build:
	go build -o commander main.go

run: build
	./commander

build-image:
	docker build -t ${IMAGE_NAME} .

run-image: build-image
	docker run ${IMAGE_NAME}
