.PHONY: all

armulator:
	go build -o bin/armulator ./cmd/armulator

testing:
	gcc -o bin/program testdata/program.c

	
all: armulator testing
