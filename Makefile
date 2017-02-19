default: all

all: build test

.SILENT: build
build:
	go build ./...

.SILENT: test
test:
	go test --short ./... --cover
