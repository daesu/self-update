BINARY_NAME=self-update
VERSION=1.1
ARCH=amd64

LINUX=$(BINARY_NAME)-linux-$(ARCH)
WINDOWS=$(BINARY_NAME)-windows-$(ARCH).exe
DARWIN=$(BINARY_NAME)-darwin-$(ARCH)

run:
	go run main.go

test:
	go test -v ./...

build: windows linux darwin

linux:
	env GOOS=linux GOARCH=$(ARCH) go build -o bin/$(LINUX) -ldflags="-s -w -X main.Version=$(VERSION)" main.go

windows:
	env GOOS=windows GOARCH=$(ARCH) go build -o bin/$(WINDOWS) -ldflags="-s -w -X main.Version=$(VERSION)" main.go

darwin:
	env GOOS=darwin GOARCH=$(ARCH) go build -o bin/$(DARWIN) -ldflags="-s -w -X main.Version=$(VERSION)" main.go

clean:
	rm -rf bin/

.PHONY: run build test