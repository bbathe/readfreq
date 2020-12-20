package := $(shell basename `pwd`)

.PHONY: default get build fmt vet

default: fmt vet

get:
	GOOS=windows GOARCH=amd64 go get -v ./...

build: default
	mkdir -p target
	rm -f target/$(package).exe
	GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w" -o target/$(package).exe

fmt:
	GOOS=windows GOARCH=amd64 go fmt ./...

vet:
	GOOS=windows GOARCH=amd64 go vet -all ./...