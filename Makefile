
GO111MODULE=off

build: fmt
	go build -o ./main/gocheck ./main

check: build
	./main/gocheck

fmt:
	gofmt -w -s *.go */*.go

init-easy:
	http_proxy=http://127.0.0.1:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/zzz.i2p
	http_proxy=http://127.0.0.1:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/echelon.i2p
	http_proxy=http://127.0.0.1:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/identiguy.i2p
	http_proxy=http://127.0.0.1:4444 curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/inr.i2p

init:
	./init.sh