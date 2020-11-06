VERSION := $(shell git describe --tags 2>/dev/null)
LDFLAGS = -X main.version=$(VERSION)

build:
	go get -u github.com/shuLhan/go-bindata/...
	$$GOPATH/bin/go-bindata -pkg dbanon -o src/bindata.go etc/*
	GO111MODULE=on go get ./...
	GO111MODULE=on go test $$GOPATH/src/github.com/mpchadwick/dbanon/src
	GO111MODULE=on go build -ldflags "$(LDFLAGS)" -o dbanon main.go
	GO111MODULE=on go test $$GOPATH/src/github.com/mpchadwick/dbanon/integration