BUILD_DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell cat VERSION)
IMAGE_NAME := prometheus-announcer
LD_FLAGS := -ldflags "-s -w -extldflags '-static' \
							-X 'main.version=$(VERSION)' \
							-X 'main.buildTime=$(BUILD_DATE)' \
							-X 'main.gitCommit=$(GIT_COMMIT)'"

protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/announce.proto

docker-build:
	docker build -f Dockerfile -t $(IMAGE_NAME):$(VERSION) .

server:
	go build -a $(LD_FLAGS) -o build/server ./cmd/server/

client:
	go build -a $(LD_FLAGS) -o build/client ./cmd/client/