build:
	go test $$HOME/go/src/github.com/mpchadwick/dbanon/src
	$$HOME/go/bin/go-bindata -pkg dbanon -o src/bindata.go etc/*
	go build -o dbanon main.go