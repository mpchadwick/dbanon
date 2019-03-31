build:
	$$HOME/go/bin/go-bindata -pkg dbanon -o src/bindata.go etc/*
	go build -o dbanon main.go