.PHONY: deps clean build test-all run-local imports vet fmt tidy lint

deps:
	go get -u ./...

clean:
	rm -rf ./cmd/*

build:
	GOOS=linux GOARCH=amd64 go build -o cmd/main main.go

test-all:
	go test -v -cover ./...


imports:
	goimports -d .

vet:
	go vet ./...

fmt:
	gofmt -d .

tidy:
	go mod tidy

lint:
	golint ./... | grep -v "comment"
