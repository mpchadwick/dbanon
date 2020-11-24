VERSION := $(shell git describe --tags 2>/dev/null)
LDFLAGS = -X main.version=$(VERSION)

build:
	rm -rf bindata
	test -z $(shell gofmt -l ./)
	go get -u github.com/shuLhan/go-bindata/...
	$$GOPATH/bin/go-bindata -pkg bindata -o bindata/bindata.go etc/*
	GO111MODULE=on go get ./...
	GO111MODULE=on go test -coverprofile=coverage.txt -covermode=atomic $$GOPATH/src/github.com/mpchadwick/dbanon/src
	GO111MODULE=on go build -ldflags "$(LDFLAGS)" -o dbanon main.go
	GO111MODULE=on go test $$GOPATH/src/github.com/mpchadwick/dbanon/integration
	rm -rf bindata

bench:
	$$GOPATH/bin/go-bindata -pkg bindata -o bindata/bindata.go etc/*
	GO111MODULE=on go test -run=XXX -bench=. -benchtime=20s $$GOPATH/src/github.com/mpchadwick/dbanon/src