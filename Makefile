.DEFAULT_GOAL := test

include test.mk

jolsat:	cmd/jolsat/jolsat.go
	go build cmd/jolsat/*.go

release: test
	go build -ldflags="-s -w"

build: jolsat

clean:
	rm jolsat
