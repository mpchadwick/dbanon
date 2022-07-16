VERSION := $(shell git describe --tags 2>/dev/null)
LDFLAGS = -X main.version=$(VERSION)

build:
	rm -rf bindata
	test -z $(shell gofmt -l ./)
	go get -u github.com/shuLhan/go-bindata/...
	go install github.com/shuLhan/go-bindata/v4/cmd/go-bindata@master
	$$GOPATH/bin/go-bindata -pkg bindata -o bindata/bindata.go etc/*
	go get ./...
	go test -coverprofile=coverage.txt -covermode=atomic -race $$GOPATH/src/github.com/mpchadwick/dbanon/src
	go build -ldflags "$(LDFLAGS)" -o dbanon main.go
	go test $$GOPATH/src/github.com/mpchadwick/dbanon/integration
	rm -rf bindata

bench:
	$$GOPATH/bin/go-bindata -pkg bindata -o bindata/bindata.go etc/*
	go test -run=XXX -bench=. -benchtime=20s $$GOPATH/src/github.com/mpchadwick/dbanon/src
	rm -rf bindata