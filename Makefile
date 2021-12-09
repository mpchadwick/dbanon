VERSION := $(shell git describe --tags 2>/dev/null)
LDFLAGS = -X main.version=$(VERSION)

build:
	test -z $(shell gofmt -l ./)
	GO111MODULE=on go get ./...
	GO111MODULE=on go test -coverprofile=coverage.txt -covermode=atomic ./src
	GO111MODULE=on go build -ldflags "$(LDFLAGS)" -o dbanon main.go
	GO111MODULE=on go test ./integration

bench:
	GO111MODULE=on go test -run=XXX -bench=. -benchtime=20s ./src
