.DEFAULT_GOAL := test

include test.mk

jolsat:	cmd/jolsat/jolsat.go
	go build $(BUILD_FLAGS) cmd/jolsat/*.go

release: BUILD_FLAGS=-ldflags="-s -w"
release: clean test

build: jolsat

clean:
	rm -f jolsat
