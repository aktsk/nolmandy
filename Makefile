NAME := nolmandy
VERSION = $(shell gobump show -r ./version)
REVISION := $(shell git rev-parse --short HEAD)

all: build

setup:
	go install golang.org/x/lint/golint
	go install golang.org/x/tools/cmd/goimports
	go install github.com/tcnksm/ghr
	go install github.com/Songmu/goxz/cmd/goxz
	go install github.com/x-motemen/gobump/cmd/gobump

test: lint
	go test ./...
	go test -race ./...

lint: setup
	golint ./...

fmt: setup
	goimports -w .

build:
	cd cmd/nolmandy; go build -o ../../bin/$(NAME)
	cd cmd/nolmandy-server; go build -o ../../bin/$(NAME)-server

clean:
	rm bin/$(NAME)

package: setup
	@sh -c "'$(CURDIR)/scripts/package.sh'"

crossbuild: setup
	goxz -pv=v${VERSION} -build-ldflags="-X main.GitCommit=${REVISION}" \
        -arch=amd64 -d=./pkg/dist/v${VERSION} \
        -n ${NAME} ./cmd/nolmandy
	goxz -pv=v${VERSION} -build-ldflags="-X main.GitCommit=${REVISION}" \
        -arch=amd64 -d=./pkg/dist/v${VERSION} \
        -n ${NAME}-server ./cmd/nolmandy-server

release: package
	ghr -u aktsk v${VERSION} ./pkg/dist/v${VERSION}

bump:
	@sh -c "'$(CURDIR)/scripts/bump.sh'"
