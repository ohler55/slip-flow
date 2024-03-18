
all: build

clean:
	rm *.so

lint:
	golangci-lint run

build:
	go build -buildmode=plugin -o flow.so *.go

test: lint
	make -C flow test

.PHONY: all build
