NAME := nolmandy

all: build

setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports

deps: setup
	glide install

test: deps lint
	go test $$(glide novendor)
	go test -race $$(glide novendor)

lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

fmt: setup
	goimports -w $$(glide nv -x)

build: deps
	go build -o bin/$(NAME) cmd/nolmandy/nolmandy.go

clean:
	rm bin/$(NAME)
