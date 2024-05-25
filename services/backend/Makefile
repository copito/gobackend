BINARY_NAME=data_quality

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

