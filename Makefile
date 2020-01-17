
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

export:
	http_proxy=http://localhost:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/export-mini.json | tee export-mini.json