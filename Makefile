# Makefile for gerrit integrator
default: build

build:
	go build

run:
	./webdav

test:
	go test

clean:
	go clean
