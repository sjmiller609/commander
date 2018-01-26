IMAGE_NAME ?= astronomerinc/ap-commander

build:
	go build -o commander main.go

build-image:
	docker build -t ${IMAGE_NAME} .
