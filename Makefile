BIN=bin

all: build

fmt:
	go fmt ./...

build: bin
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o $(BIN)/fswatcher

install:
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go install

clean:
	go clean -i
	rm -rf $(BIN)

$(BIN):
	mkdir -p $(BIN)

.PHONY: fmt install clean all

