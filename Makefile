.PHONY: sort generator validator

all: proto build

proto: pkg/mafiapb/mafia.proto
	@echo "Generating proto objects..."
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pkg/mafiapb/mafia.proto

build:
	go build -o bin/server cmd/server/*
	go build -o bin/client cmd/client/*
