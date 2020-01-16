
GO111MODULE=on

build: fmt
	go build -o ./main/gocheck ./main

check: build
	./main/gocheck

fmt:
	gofmt -w -s *.go */*.go

force-check:
	./init.sh