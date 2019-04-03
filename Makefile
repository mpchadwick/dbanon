build:
	go test $$GOPATH/src/github.com/mpchadwick/dbanon/src
	$$GOPATH/bin/go-bindata -pkg dbanon -o src/bindata.go etc/*
	go build -o dbanon main.go