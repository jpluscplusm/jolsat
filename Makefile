.DEFAULT_GOAL := test

include test.mk

jolsat:	jolsat.go
	go build

release: test
	go build -ldflags="-s -w"

build: jolsat

clean:
	rm jolsat
