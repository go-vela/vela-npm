# SPDX-License-Identifier: Apache-2.0

# capture the current date we build the application from
BUILD_DATE = $(shell date +%Y-%m-%dT%H:%M:%SZ)

# check if a git commit sha is already set
ifndef GITHUB_SHA
	# capture the current git commit sha we build the application from
	GITHUB_SHA = $(shell git rev-parse HEAD)
endif

# check if a git tag is already set
ifndef GITHUB_TAG
	# capture the current git tag we build the application from
	GITHUB_TAG = $(shell git describe --tag --abbrev=0)
endif

# check if a go version is already set
ifndef GOLANG_VERSION
	# capture the current go version we build the application from
	GOLANG_VERSION = $(shell go version | awk '{ print $$3 }')
endif

# create a list of linker flags for building the golang application
LD_FLAGS = -X github.com/go-vela/vela-npm/version.Commit=${GITHUB_SHA} -X github.com/go-vela/vela-npm/version.Date=${BUILD_DATE} -X github.com/go-vela/vela-npm/version.Go=${GOLANG_VERSION} -X github.com/go-vela/vela-npm/version.Tag=${GITHUB_TAG}

# The `clean` target is intended to clean the workspace
# and prepare the local changes for submission.
#
# Usage: `make clean`
.PHONY: clean
clean: tidy vet fmt fix reset

# The `run` target is intended to build and
# execute the Docker image for the plugin.
#
# Usage: `make run`
.PHONY: run
run: build docker-build docker-run

# The `tidy` target is intended to clean up
# the Go module files (go.mod & go.sum).
#
# Usage: `make tidy`
.PHONY: tidy
tidy:
	@echo
	@echo "### Tidying Go module"
	@go mod tidy

# The `vet` target is intended to inspect the
# Go source code for potential issues.
#
# Usage: `make vet`
.PHONY: vet
vet:
	@echo
	@echo "### Vetting Go code"
	@go vet ./...

# The `fmt` target is intended to format the
# Go source code to meet the language standards.
#
# Usage: `make fmt`
.PHONY: fmt
fmt:
	@echo
	@echo "### Formatting Go Code"
	@go fmt ./...

# The `fix` target is intended to rewrite the
# Go source code using old APIs.
#
# Usage: `make fix`
.PHONY: fix
fix:
	@echo
	@echo "### Fixing Go Code"
	@go fix ./...

# The `reset` target is intended to clean up
# env files and tear down Verdaccio
#
# Usage: `make fix`
.PHONY: reset
reset:
	@echo
	@echo "### Resetting local env"
	@rm -f .env .env-docker
	@touch .env-docker .env
	@docker compose down

# The `test` target is intended to run
# the tests for the Go source code.
#
# Usage: `make test`
.PHONY: test
test:
	@echo
	@echo "### Testing Go Code"
	@go test -race ./...

# The `test-cover` target is intended to run
# the tests for the Go source code and then
# open the test coverage report.
#
# Usage: `make test-cover`
.PHONY: test-cover
test-cover:
	@echo
	@echo "### Creating test coverage report"
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo
	@echo "### Opening test coverage report"
	@go tool cover -html=coverage.out

# The `build` target is intended to compile
# the Go source code into a binary.
#
# Usage: `make build`
.PHONY: build
build:
	@echo
	@echo "### Building release/vela-npm binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-ldflags '${LD_FLAGS}' \
		-o release/vela-npm \
		github.com/go-vela/vela-npm/cmd/vela-npm

# The `build-static` target is intended to compile
# the Go source code into a statically linked binary.
#
# Usage: `make build-static`
.PHONY: build-static
build-static:
	@echo
	@echo "### Building static release/vela-npm binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-ldflags '-s -w -extldflags "-static" ${LD_FLAGS}' \
		-o release/vela-npm \
		github.com/go-vela/vela-npm/cmd/vela-npm

# The `build-static-ci` target is intended to compile
# the Go source code into a statically linked binary
# when used within a CI environment.
#
# Usage: `make build-static-ci`
.PHONY: build-static-ci
build-static-ci:
	@echo
	@echo "### Building CI static release/vela-npm binary"
	@go build -a \
		-ldflags '-s -w -extldflags "-static" ${LD_FLAGS}' \
		-o release/vela-npm \
		github.com/go-vela/vela-npm/cmd/vela-npm

# The `check` target is intended to output all
# dependencies from the Go module that need updates.
#
# Usage: `make check`
.PHONY: check
check: check-install
	@echo
	@echo "### Checking dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -update

# The `check-direct` target is intended to output direct
# dependencies from the Go module that need updates.
#
# Usage: `make check-direct`
.PHONY: check-direct
check-direct: check-install
	@echo
	@echo "### Checking direct dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -direct

# The `check-full` target is intended to output
# all dependencies from the Go module.
#
# Usage: `make check-full`
.PHONY: check-full
check-full: check-install
	@echo
	@echo "### Checking all dependencies for updates"
	@go list -u -m -json all | go-mod-outdated

# The `check-install` target is intended to download
# the tool used to check dependencies from the Go module.
#
# Usage: `make check-install`
.PHONY: check-install
check-install:
	@echo
	@echo "### Installing psampaz/go-mod-outdated"
	@go get -u github.com/psampaz/go-mod-outdated

# The `bump-deps` target is intended to upgrade
# non-test dependencies for the Go module.
#
# Usage: `make bump-deps`
.PHONY: bump-deps
bump-deps: check
	@echo
	@echo "### Upgrading dependencies"
	@go get -u ./...

# The `bump-deps-full` target is intended to upgrade
# all dependencies for the Go module.
#
# Usage: `make bump-deps-full`
.PHONY: bump-deps-full
bump-deps-full: check
	@echo
	@echo "### Upgrading all dependencies"
	@go get -t -u ./...

# The `docker-build` target is intended to build
# the Docker image for the plugin.
#
# Usage: `make docker-build`
.PHONY: docker-build
docker-build:
	@echo
	@echo "### Building vela-npm:local image"
	@docker build --no-cache -t vela-npm:local .

# The `docker-run` target is intended to execute
# the Docker image for the plugin.
#
# Usage: `make docker-run`
VERDACCIO_TOKEN := $(shell cat .env)
PROJECT_PATH ?= $(shell pwd)/examples/example-0
.PHONY: docker-run
docker-run:
	@echo
	@echo "### Executing vela-npm:local image"
	docker run --rm \
		-e PARAMETER_USERNAME=testuser \
		-e NPM_PASSWORD=testpass \
		-e NPM_TOKEN=$(VERDACCIO_TOKEN) \
		-e PARAMETER_EMAIL= \
		-e PARAMETER_REGISTRY="http://localhost:4873" \
		-e PARAMETER_SKIP_PING= \
		-e PARAMETER_STRICT_SSL=false \
		-e PARAMETER_ALWAYS_AUTH= \
		-e PARAMETER_DRY_RUN= \
		-e PARAMETER_LOG_LEVEL=trace \
		-e PARAMETER_AUDIT_LEVEL=critical \
		-e PARAMETER_TAG= \
		-e PARAMETER_ACCESS=public \
		-e PARAMETER_WORKSPACES= \
		-e PARAMETER_WORKSPACE= \
		-v $(PROJECT_PATH):/workspace \
		-w /workspace \
		--network host \
		vela-npm:local

# The `test-e2e` target is intended to execute
# an end-to-end test with Verdaccio with all
# setup and teardown needed
#
# Usage: `make test-e2e`
.PHONY: test-e2e
test-e2e: clean prepare-examples user docker-run

# The `docker-example` target is intended to execute
# an end-to-end test with Verdaccio
#
# Usage: `make docker-example`
.PHONY: docker-example
docker-example:
	@echo
	@echo "### Testing plugin against Verdaccio"
	@docker-compose up --abort-on-container-exit

# The `prepare-examples` target is intended to execute
# fresh installs on all example npm packages
#
# Usage: `make prepare-examples`
.PHONY: prepare-examples
prepare-examples:
	@echo
	@echo "### Cleaning and preparing example packages"
	@cd examples \
		&& rm -rf node_modules package-lock.json \
		&& npm i -ws \
		&& cd example-0 \
		&& rm -rf node_modules package-lock.json \
		&& npm i \
		&& cd ..
	@echo "Examples ready!"

# The `user` target is stand up a Verdaccio instance
# and create a user (testuser:testpass), storing the
# generated user's token in .env file
#
# Usage: `make user`
.PHONY: user
user:
	scripts/make-user.sh

# The `registry` target is intended to stand
# up a Verdaccio instance
#
# Usage: `make registry`
.PHONY: registry
registry:
	@scripts/open-registry.sh
