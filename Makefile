.PHONY: all build binrun run fmt clean docker

all: build

build:
	go build -o ./build/app ./cmd/app

binrun: build
	./build/app

run:
	go run main.go

fmt:
	go mod tidy
	go fmt ./...

clean:
	rm ./build/*

docker:
	docker build -t server . --build-arg DEV=true
