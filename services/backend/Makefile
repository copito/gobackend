BINARY_NAME=data_quality

PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := protoc-gen-grpc-gateway

# point to proto files outside service
PROTO_DIR := ./../idl/grpc
OUT_DIR := ./../idl-gen

# find all proto files
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

build:
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin src/main.go
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux src/main.go
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows src/main.go

run: build
	bin/${BINARY_NAME}

clean:
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux
	rm bin/${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all


list_installed_packages:
	go list -f "{{.ImportPath}} {{.Imports}}" ./...


## Minikube commands
# this shows the eval command to point to minikube's docker registry
minikube_docker:
	minikube docker-env

minikube_build:
	docker build -t local/devex:v1 ./deploy/Dockerfile && docker images

