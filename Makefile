PACKAGES := github.com/greendinosaur/gh-commit-info
WORKDIR := $(PACKAGES)/src/api

all: build silent-test

build:
	go build -v ./...
	cd src/api;go build -o main .

test:
	go test -v -coverprofile cp.out -count=1 -parallel=1 $(PACKAGES)/...
	go tool cover -html=cp.out -o coverage.html

silent-test:
	go test -coverprofile cp.out -count=1 -parallel=1 $(PACKAGES)/...
	go tool cover -html=cp.out -o coverage.html
	
format:
	go fmt $(PACKAGES)/...

clean:
	go clean
