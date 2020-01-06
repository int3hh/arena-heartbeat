-include .env

build:
	@echo " > Building binary ... "
	go build
clean:
	go clean