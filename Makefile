
GO111MODULE=off

build: fmt
	go build -o ./main/gocheck ./main

check: build
	./main/gocheck

fmt:
	gofmt -w -s *.go */*.go