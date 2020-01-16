
GO111MODULE=on

build: fmt
	go build -o ./main/gocheck ./main

check: build
	./main/gocheck

fmt:
	gofmt -w -s *.go */*.go
	prettier --write script.js

force-check:
	./init.sh