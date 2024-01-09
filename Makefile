DOCKER_USERNAME ?= brownbarg
APPLICATION_NAME ?= craise
APPLICATION_DIR := ./build
APPLICATION_FILEPATH := ./build/craise
HELMPACKAGE := "craise"

.PHONY: build test clean run docker-run docker-build docker-push docker-tag helm-create helm-package helm-test kdeploy

test:
	go test ./... -v

build:
	rm -rf build
	@go build -v --o $(APPLICATION_FILEPATH) main.go

build-debug:
	rm -rf build
	@go build -v --o $(APPLICATION_FILEPATH) -ldflags "-X main.Debug=true" main.go

build-debug-ff:
	rm -rf build
	@#go build -v --o $(APPLICATION_FILEPATH) -ldflags "-X main.Debug=true -X main.FeatureFlags=GracefulTermination,PlayGame" main.go
	@go build -v --o $(APPLICATION_FILEPATH) -ldflags "-X main.Debug=true -X main.FeatureFlags=PlayGame" main.go

clean:
	@rm -rf ./build

run: build
	@echo "Running with args: $(ARGS)"
	@$(APPLICATION_FILEPATH) $(ARGS)

run-debug: build-debug
	@echo "Running with args: $(ARGS)"
	@$(APPLICATION_FILEPATH) $(ARGS)

run-debug-ff: build-debug-ff
	@echo "Running with args: $(ARGS)"
	@$(APPLICATION_FILEPATH) $(ARGS)