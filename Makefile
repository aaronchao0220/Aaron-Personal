VERSION ?= latest
SERVICE_NAME := usage-telemetry-publisher
DOCKER_IMAGE := qlik/$(SERVICE_NAME)
ARGS ?= ""

BUILD_TIME = `date -u +%Y-%m-%dT%H:%M:%SZ`
REVISION ?= `git rev-parse HEAD 2>/dev/null`
BRANCH ?= `git rev-parse --abbrev-ref HEAD`

BRANCH_FLAG=-X `go list ./cmd/version`.Branch=${BRANCH}
REVISION_FLAG=-X `go list ./cmd/version`.Revision=${REVISION}
VERSION_FLAG=-X `go list ./cmd/version`.Version=${VERSION}
BUILD_TIME_FLAG=-X `go list ./cmd/version`.BuildTime=${BUILD_TIME}

# Rebuild dependencies
mod:
	@go mod tidy

# Update dependencies
mod-update:
	@go get -u ./...
	@$(MAKE) mod

# Lint the code
lint:
	./scripts/lint.sh

# Compile Go packages and dependencies
build:
	@GOOS=$(shell echo $* | cut -f1 -d-) GOARCH=$(shell echo $* | cut -f2 -d- | cut -f1 -d.) go build -o ./"$(SERVICE_NAME)" \
	-ldflags "$(REVISION_FLAG) $(VERSION_FLAG) $(BUILD_TIME_FLAG)" ./cmd/main/main.go

build-docker-image:
	export DOCKER_BUILDKIT=1 && docker build --platform linux/amd64 --tag $(DOCKER_IMAGE)$(IMAGE_NAME_SUFFIX):$(VERSION) --file ./docker/dockerfile --target $(BUILD_TARGET) \
	--build-arg CREATED=$(BUILD_TIME) \
	--build-arg REVISION=$(REVISION) \
	--build-arg VERSION=$(VERSION) \
	--ssh default \
	.
	docker tag $(DOCKER_IMAGE)$(IMAGE_NAME_SUFFIX):$(VERSION) $(DOCKER_IMAGE)$(IMAGE_NAME_SUFFIX):latest

# Build the Docker image
build-docker:
	@IMAGE_NAME_SUFFIX= BUILD_TARGET=production $(MAKE) build-docker-image

# Build the Docker test image
build-test-docker:
	@IMAGE_NAME_SUFFIX=-test BUILD_TARGET=test $(MAKE) build-docker-image

# Run unit tests
test-unit:
	./scripts/test-unit.sh

# Upload coverage
upload-coverage:
	./scripts/report-coverage.sh

test-component: start-dependencies
	./scripts/test-component.sh


# Run component tests locally against dependencies in docker containers spun up by the tests
test-component-docker: build-docker build-test-docker test-component-nobuild

test-component-nobuild: start-nobuild
	docker compose -f docker/docker-compose.yaml up test --exit-code-from test

test-e2e-docker: build-docker build-test-docker test-e2e-nobuild

test-e2e-nobuild: start-nobuild
	docker compose -f docker/docker-compose.yaml up test_e2e --exit-code-from test_e2e

start-nobuild: start-dependencies
	COMPOSE_IGNORE_ORPHANS=1 docker compose -f ./docker/docker-compose.yaml up -d $(SERVICE_NAME)

start-dependencies: stop
	COMPOSE_IGNORE_ORPHANS=1 docker compose -f ./docker/docker-compose.yaml up -d solace server-mocks ldRelay

stop:
	docker compose -f ./docker/docker-compose.yaml down -v --remove-orphans
