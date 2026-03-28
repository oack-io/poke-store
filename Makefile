.PHONY: build build-server build-web install-web run dev test test-e2e test-oack lint clean docker docker-push docker-run

APP_NAME := poke-store
BUILD_DIR := bin
VERSION ?= dev
COMMIT_SHA := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_MSG := $(shell git log -1 --format='%s')
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commitSHA=$(COMMIT_SHA) -X main.buildTime=$(BUILD_TIME)"
DOCKER_REPO ?= oack/poke-store
DOCKER_PLATFORMS ?= linux/amd64,linux/arm64

build: build-web sync-static build-server

build-server:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

build-web:
	cd web && npm run build

sync-static:
	rm -rf cmd/server/static
	cp -r web/dist cmd/server/static

install-web:
	cd web && npm install

run: build
	$(BUILD_DIR)/$(APP_NAME)

dev:
	go run $(LDFLAGS) ./cmd/server

test:
	go test ./...

test-e2e:
	cd web && npx playwright test

test-oack:
	cd oack-checks && BASE_URL=$(or $(BASE_URL),http://localhost:6001) npx playwright test

lint:
	golangci-lint run --tests=false ./...

lint-web:
	cd web && npm run lint

docker:
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(DOCKER_REPO):$(VERSION) \
		-t $(DOCKER_REPO):latest .

docker-push:
	docker buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--build-arg VERSION=$(VERSION) \
		-t $(DOCKER_REPO):$(VERSION) \
		-t $(DOCKER_REPO):latest \
		--push .

docker-run:
	docker run --rm -p 6001:6001 \
		-e SESSION_SECRET=$(or $(SESSION_SECRET),change-me-in-prod) \
		$(DOCKER_REPO):latest

clean:
	rm -rf $(BUILD_DIR) web/dist
