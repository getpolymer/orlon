# Makefile for Orlon
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_DIR=bin
BINARY_NAME=orlon
ORG_PATH="github.com/getpolymer"
PROJ="orlon"
REPO_PATH="${ORG_PATH}/${PROJ}"

all: build

build:
	$(GOBUILD) -i -o "$(BINARY_DIR)/$(BINARY_NAME)" "$(REPO_PATH)/cmd/orlon"

clean:
	$(GOCLEAN)
	rm -rf "$(BINARY_DIR)"
