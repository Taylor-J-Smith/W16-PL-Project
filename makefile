build: path
	go build server.go
	./server

run: path
	go run server.go

path:
	export PATH=$PATH:/usr/local/go/bin
	export GOPATH=$HOME/Go-Code
