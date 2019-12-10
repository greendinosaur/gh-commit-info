PACKAGES := github.com/greendinosaur/gh-commit-info
DEPENDENCIES := 

all: build silent-test

build:
	go build $(PACKAGES)/...

test:
	go test -v -coverprofile cp.out -count=1 -parallel=1 $(PACKAGES)/...


silent-test:
	go test -coverprofile cp.out -count=1 -parallel=1 $(PACKAGES)/...

format:
	go fmt $(PACKAGES)/...

clean:
	go clean