SHELL = /bin/sh
.DEFAULT_GOAL := buildAndRun

# VARIABLES
BIN = ./bin

$(BIN)/:
	mkdir -p $@

.PHONE:clean
clean: $(BIN)/

build: clean
	go build -o ./bin cmd/app/app.go
	cp config.yaml ./bin/config.yaml

run:
	./bin/app

buildAndRun: build run