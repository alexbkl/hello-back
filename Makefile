# Copyright © 2023 Hello Storage Inc. All rights reserved.
#
export GO111MODULE=on

# 
APP_NAME=hello
DOCKER_COMPOSE=docker compose
GOTEST=go test

# 
develop:
	scripts/docker/build.sh develop
developx:
	scripts/docker/buildx.sh develop
start-local:
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml up -d --wait
stop-local:
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml stop
logs:
	$(DOCKER_COMPOSE) logs -f
build-go:
	rm -f build/$(APP_NAME)
	scripts/build.sh debug $(APP_NAME)
test:
	$(info Running all Go tests...)
	$(GOTEST) -parallel 1 -count 1 -cpu 1 -tags slow -timeout 20m ./pkg/... ./internal/...
fmt:
	go fmt ./pkg/... ./internal/... ./cmd/...
	gofmt -w -s pkg internal cmd
	goimports -w pkg internal cmd
install:
	go get .
tidy:
	go mod tidy

# TO-DO drone configure
# drone:


.PHONY: all test clean