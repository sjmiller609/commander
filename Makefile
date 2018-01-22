# IMAGE_NAME ?= astronomerinc/commander

build:
	go build -o commander main.go

run: build
	./commander
