# Version control options
GIT 	 := git
VERSION  := $(shell $(GIT) describe --match 'v[0-9]*' --dirty='.m' --always --tags)
REVISION := $(shell $(GIT) rev-parse HEAD)$(shell if ! $(GIT) diff --no-ext-diff --quiet --exit-code; then echo .m; fi)

# Golang options
GO       ?= go
BINDATA  ?= bindata
MOCKGEN  ?= mockgen
PKG      := github.com/nomkhonwaan/myblog
TAGS     :=
LDFLAGS  :=
GOFLAGS  :=
BINDIR   := $(CURDIR)/bin

# Node.js options
NPM      ?= npm
NPX      ?= npx
NG       := $(NPX) ng
WEBDIR   := $(CURDIR)/web

# Docker options
DOCKER   := docker

.PHONY: all
all: clean install build install-web build-web
	
.PHONY: install
install:
	$(GO) mod download

.PHONY: install-web
install-web:
	cd web && $(NPM) install --silent

.PHONY: clean
clean:
	rm -rf $(BINDIR)/myblog && \
	rm -rf $(CURDIR)/coverage.out && \
	rm -rf $(CURDIR)/vendor && \
	rm -rf $(WEBDIR)/dist && \
	rm -rf $(WEBDIR)/node_modules

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: test
test:
	$(GO) test -v ./... -race -coverprofile=coverage.out -covermode=atomic
	
.PHONY: bindata
bindata:
	$(BINDATA) -o ./pkg/data/data.go ./data/...

.PHONY: build
build:
	$(GO) build $(GOFLAGS) -tags '$(TAGS)' -ldflags '-X $(PKG)/cmd/myblog.Version=$(VERSION) -X $(PKG)/cmd/myblog.Revision=$(REVISION) $(LDFLAGS)' -o $(BINDIR)/myblog main.go

.PHONY: build-web
build-web:
	mv $(WEBDIR)/src/environments/environment.prod.ts $(WEBDIR)/src/environments/environment.prod.original.ts && \
	VERSION=$(VERSION) REVISION=$(REVISION) envsubst < $(WEBDIR)/src/environments/environment.prod.original.ts > $(WEBDIR)/src/environments/environment.prod.ts && \
	cd $(WEBDIR) && \
	$(NG) build --prod --buildOptimizer --vendorChunk --source-map=false && \
	rm -f $(WEBDIR)/src/environments/environment.prod.ts && \
	mv $(WEBDIR)/src/environments/environment.prod.original.ts $(WEBDIR)/src/environments/environment.prod.ts

.PHONY: build-docker
build-docker:
	$(DOCKER) build --file build/package/Dockerfile --tag nomkhonwaan/myblog:latest .

.PHONY: build-docker-all-in-one
build-docker-all-in-one:
	$(DOCKER) build --build-arg NPM_AUTH_TOKEN=${NPM_AUTH_TOKEN} --file build/package/all-in-one/Dockerfile --tag nomkhonwaan/myblog-all-in-one:latest .

.PHONY: build-docker-all-in-one-ci
build-docker-all-in-one-ci:
	$(DOCKER) build --file build/package/all-in-one-ci/Dockerfile --tag nomkhonwaan/myblog-all-in-one:latest .
