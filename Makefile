APP_NAME = app
DOCKER_REPO = ""
RELEASE_VERSION := ""

run:
	@go mod tidy
	@go run -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" main.go

build:
	go build -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o ./bin/$(APP_NAME) main.go

window:
	echo "Compiling for windows"
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/window/$(APP_NAME).exe main.go
	echo "Output stored in ./bin"

linux:
	echo "Compiling for linux"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/linux/$(APP_NAME) main.go
	echo "Output stored in ./bin"

alpine:
	echo "Compiling for docker alpine"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -a -installsuffix cgo -o bin/alpine/$(APP_NAME) main.go
	echo "Output stored in ./bin"

docker:
	docker build -t $(APP_NAME) --platform linux/amd64 --build-arg APP_NAME=$(APP_NAME) --build-arg RELEASE_VERSION=$(RELEASE_VERSION) .

hub:
	docker tag $(APP_NAME) $(DOCKER_REPO):$(RELEASE_VERSION)
	docker push $(DOCKER_REPO):$(RELEASE_VERSION)

orm:
	sqlboiler psql --add-global-variants

grpc:
	protoc --go_out=./pkg/core/pb --go_opt=paths=source_relative   --go-grpc_out=./pkg/core/pb --go-grpc_opt=paths=source_relative   --proto_path=absolute_path_to_proto_directory absolute_path_to_proto_directory/*.proto
