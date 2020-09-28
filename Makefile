SHELL    = /bin/bash
AUTHOR   = eirsyl
PACKAGE  = feedy

GIT_USER       = eirsyl
GIT_REPOSITORY = feedy
REGISTRY_OWNER = eirsyl
REGISTRY_IMAGE = feedy

DATE    ?= $(shell date +%FT%T%z)
VERSION  = $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo v0)
BIN      = $(GOPATH)/bin
IMPORT   = github.com/$(AUTHOR)/$(PACKAGE)
BASE     = $(shell pwd)
PKGS     = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -Ev "vendor"))
TESTPKGS = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))

GO      = go
GODOC   = godoc
GOFMT   = gofmt
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

export PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all
all: test build | $(BASE) ;
	$Q

$(BASE): ; $(info $(M) checking GOPATH…)
	@echo

# Tools
.PHONY: gopath
gopath:
	@echo $(GOPATH)

GOLANGCI-LINT = $(BIN)/golangci-lint
$(BIN)/golangci-lint: | $(BASE) ; $(info $(M) building golangci-lint…)
	$Q go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0

GOCOVMERGE = $(BIN)/gocovmerge
$(BIN)/gocovmerge: | $(BASE) ; $(info $(M) building gocovmerge…)
	$Q go get -u github.com/wadey/gocovmerge

GOCOV = $(BIN)/gocov
$(BIN)/gocov: | $(BASE) ; $(info $(M) building gocov…)
	$Q go get -u github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: | $(BASE) ; $(info $(M) building gocov-xml…)
	$Q go get -u github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: | $(BASE) ; $(info $(M) building go2xunit…)
	$Q go get -u github.com/tebeka/go2xunit

STRINGER = $(BIN)/stringer
$(BIN)/stringer: | $(BASE) ; $(info $(M) building stringer…)
	$Q go get -u golang.org/x/tools/cmd/stringer

GITHUBRELEASE = $(BIN)/github-release
$(BIN)/github-release: | $(BASE) ; $(info $(M) building github-release…)
	$Q go get -u github.com/aktau/github-release

TEST-RESULTS = $(BASE)/test-results
$(BASE)/test-results: | $(BASE) ; $(info $(M) creating test-results…)
	$Q mkdir $(BASE)/test-results

ARTIFACTS = $(BASE)/artifacts
$(BASE)/artifacts: | $(BASE) ; $(info $(M) creating artifacts…)
	$Q mkdir $(BASE)/artifacts

MIGRATE = $(BIN)/migrate
$(BIN)/migrate: | $(BASE) ; $(info $(M) building migrate…)
	$Q go get -u github.com/golang-migrate/migrate/v4/cmd/migrate

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml test
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
test: generate fmt lint | $(BASE) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q cd $(BASE) && $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: generate fmt lint | $(BASE) $(GO2XUNIT) $(TEST-RESULTS) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q cd $(BASE) && 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee $(TEST-RESULTS)/tests.output
	$(GO2XUNIT) -fail -input $(TEST-RESULTS)/tests.output -output $(TEST-RESULTS)/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML) $(ARTIFACTS)
test-coverage: COVERAGE_DIR := $(ARTIFACTS)/coverage
test-coverage: generate test-coverage-tools | $(BASE) ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)
	$Q cd $(BASE) && for pkg in $(TESTPKGS); do \
		$(GO) test \
			-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
					grep '^$(PACKAGE)/' | grep -Ev 'vendor/' | \
					tr '\n' ',')$$pkg \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_DIR)/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	 done
	$Q $(GOCOVMERGE) $(COVERAGE_DIR)/*.cover > $(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: lint
lint: generate | $(BASE) $(GOLANGCI-LINT) ; $(info $(M) running golangci-lint…) @ ## Run golangci-lint
	$Q cd $(BASE) && $(GOLANGCI-LINT) run ./...

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./...); do \
		$(GOFMT) -s -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

# Build

generate: | $(BASE) $(STRINGER) ; $(info $(M) generating code…) @ ## Generate code
	$Q cd $(BASE) && $(GO) generate ./...

build: generate | $(BASE) ; $(info $(M) building executable…) @ ## Build program binary
	$Q cd $(BASE) && CGO_ENABLED=0 GO111MODULE=on $(GO) build \
		-tags release \
		-ldflags '-X $(IMPORT)/internal.Version="$(if $(OUT_VERSION),$(OUT_VERSION),$(VERSION))" -X $(IMPORT)/internal.BuildDate=$(DATE)' \
		-o $(if $(OUT),$(OUT),bin/$(PACKAGE)) main.go

# Docker

DOCKER_IMAGE=$(REGISTRY_OWNER)/$(REGISTRY_IMAGE)
DOCKER_REGISTRY_IMAGE=$(REGISTRY_OWNER)/$(REGISTRY_IMAGE)

.PHONY: docker-container
docker-container: | $(BASE) ; $(info $(M) building container…) @ ## Build container
	$Q docker build --pull -t $(DOCKER_IMAGE):$(VERSION) .
	$Q $(MAKE) docker-tag-helper

.PHONY: docker-push
docker-push: docker-container | $(BASE) ; $(info $(M) pushing container…) @ ## Push container
	$Q $(MAKE) docker-push-helper

# Docker helpers

.PHONY: docker-tag-helper
docker-tag-helper: | $(BASE) ; $(info $(M) tagging docker container…) @
	$Q docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_REGISTRY_IMAGE):latest
	$Q docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_REGISTRY_IMAGE):$(VERSION)

.PHONY: docker-push-helper
docker-push-helper: | $(BASE) ; $(info $(M) tagging docker container…) @
	$Q docker push $(DOCKER_IMAGE):$(VERSION) $(DOCKER_REGISTRY_IMAGE):latest
	$Q docker push $(DOCKER_IMAGE):$(VERSION) $(DOCKER_REGISTRY_IMAGE):$(VERSION)

# Release

# arm
bin/linux/arm/5/$(PACKAGE): | $(BASE) ; $(info $(M) building arm-5 executable…) @
	$Q GOARM=5 GOARCH=arm GOOS=linux OUT="$@" $(MAKE) build
bin/linux/arm/7/$(PACKAGE): | $(BASE) ; $(info $(M) building arm-7 executable…) @
	$Q GOARM=7 GOARCH=arm GOOS=linux OUT="$@" $(MAKE) build
bin/linux/arm64/$(PACKAGE): | $(BASE) ; $(info $(M) building arm64 executable…) @
	$Q GOARCH=arm64 GOOS=linux OUT="$@" $(MAKE) build

# 386
bin/darwin/386/$(PACKAGE): | $(BASE) ; $(info $(M) building darwin-386 executable…) @
	$Q GOARCH=386 GOOS=darwin OUT="$@" $(MAKE) build
bin/linux/386/$(PACKAGE): | $(BASE) ; $(info $(M) building linux-386 executable…) @
	$Q GOARCH=386 GOOS=linux OUT="$@"$(MAKE) build
bin/windows/386/$(PACKAGE): | $(BASE) ; $(info $(M) building windows-386 executable…) @
	$Q GOARCH=386 GOOS=windows OUT="$@" $(MAKE) build

# amd64
bin/freebsd/amd64/$(PACKAGE): | $(BASE) ; $(info $(M) building freebsd-amd64 executable…) @
	$Q GOARCH=amd64 GOOS=freebsd OUT="$@" $(MAKE) build
bin/darwin/amd64/$(PACKAGE): | $(BASE) ; $(info $(M) building darwin-amd64 executable…) @
	$Q GOARCH=amd64 GOOS=darwin OUT="$@" $(MAKE) build
bin/linux/amd64/$(PACKAGE): | $(BASE) ; $(info $(M) building linux-amd64 executable…) @
	$Q GOARCH=amd64 GOOS=linux OUT="$@" $(MAKE) build
bin/windows/amd64/$(PACKAGE).exe: | $(BASE) ; $(info $(M) building windows-amd64 executable…) @
	$Q GOARCH=amd64 GOOS=windows OUT="$@" $(MAKE) build

UNIX_EXECUTABLES := \
	linux/arm/5/$(PACKAGE) \
	linux/arm/7/$(PACKAGE) \
	linux/arm64/$(PACKAGE) \
	darwin/amd64/$(PACKAGE) \
	freebsd/amd64/$(PACKAGE) \
	linux/amd64/$(PACKAGE)

COMPRESSED_EXECUTABLES=$(UNIX_EXECUTABLES:%=%.bz2)
COMPRESSED_EXECUTABLE_TARGETS=$(COMPRESSED_EXECUTABLES:%=bin/%)

# UPLOAD_CMD = $(GITHUBRELEASE) upload -u $(GIT_USER) -r $(GIT_REPOSITORY) -t $(VERSION) -n $(subst /,-,$(FILE)) -f bin/$(FILE)
UPLOAD_CMD = echo "upload -u $(GIT_USER) -r $(GIT_REPOSITORY) -t $(VERSION) -n $(subst /,-,$(FILE)) -f bin/$(FILE)"

%.bz2: %
	$Q bzip2 -c < "$<" > "$@"

.PHONY: release
release: clean | $(BASE) $(GITHUBRELEASE) ; $(info $(M) releasing application…) @ ## Upload release to GitHub
	$Q OUT_VERSION=$(VERSION) $(MAKE) $(COMPRESSED_EXECUTABLE_TARGETS)
	$Q git log --format=%B $(VERSION) -1 | \
		$(GITHUBRELEASE) release -u $(GIT_USER) -r $(GIT_REPOSITORY) \
		-t $(VERSION) -n $(VERSION) -d - || true
	$Q $(foreach FILE,$(COMPRESSED_EXECUTABLES),$(UPLOAD_CMD);)

# Database

.PHONY: database-create-migration
database-create-migration: | $(BASE) $(MIGRATE) ; $(info $(M) generating migration file…) @ ## Generate migration file
	$Q cd $(BASE) && $(MIGRATE) create -ext sql -dir contrib/migrations/server -seq $$MIGRATION

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin vendor
	@rm -rf $(TEST-RESULTS)
	@rm -rf $(ARTIFACTS)

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
