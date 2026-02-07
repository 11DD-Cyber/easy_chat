BUILD_DIR := apps/user/rpc
BINARY := build/user_rpc
DOCKERFILE := deploy/dockerfile/Dockerfile_user_rpc_test
IMAGE := easy-chat/user-rpc:latest

.PHONY: build docker-build run-local clean

build:
    mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY) $(BUILD_DIR)

docker-build: build
	docker build -f $(DOCKERFILE) -t $(IMAGE) .

run-local: build
	./$(BINARY) -f apps/user/rpc/etc/user.yaml

clean:
	rm -f $(BINARY)