NAME := nolmandy

all: build

setup:
	go get golang.org/x/vgo
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports

test: lint
	vgo test ./...
	vgo test -race ./...

lint: setup
	vgo vet ./...
	golint ./...

fmt: setup
	goimports -w .

build:
	vgo build -o bin/$(NAME) cmd/nolmandy/nolmandy.go
	vgo build -o bin/$(NAME)-server cmd/nolmandy/nolmandy_server.go

clean:
	rm bin/$(NAME)

