# Version control options
GIT 	 := git
VERSION  := $(shell $(GIT) describe --match 'v[0-9]*' --dirty='.m' --always)
REVISION := $(shell $(GIT) rev-parse HEAD)$(shell if ! $(GIT) diff --no-ext-diff --quiet --exit-code; then echo .m; fi)

# Golang options
GO       ?= go
DEP	 := dep
PKG	 := github.com/nomkhonwaan/myblog
TAGS     :=
LDFLAGS  :=
GOFLAGS  :=
BINDIR   := $(CURDIR)/bin

# Docker options
DOCKER   := docker

.PHONY: all
all: clean install build


.PHONY: clean
clean:
	@rm -rf $(BINDIR)
 
.PHONY: install
install:
	$(DEP) ensure

.PHONY: test
test:
	$(GO) test ./...
	
.PHONY: build
build:
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '-X main.version=$(VERSION) -X main.revision=$(REVISION) $(LDFLAGS)' $(PKG)/cmd/myblog

.PHONY: docker-build
docker-build:
	docker build -f build/package/Dockerfile -t nomkhonwaan/myblog:latest .
