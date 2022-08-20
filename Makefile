.PHONY: all

armulator:
	go build -o bin/armulator ./cmd/armulator

testing:
	aarch64-linux-gnu-gcc -o bin/program testdata/program.c
	
all: armulator testing
