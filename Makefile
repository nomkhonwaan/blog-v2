# Version control options
GIT 	 := git
VERSION  := $(shell $(GIT) describe --match 'v[0-9]*' --dirty='.m' --always)
REVISION := $(shell $(GIT) rev-parse HEAD)$(shell if ! $(GIT) diff --no-ext-diff --quiet --exit-code; then echo .m; fi)

# Golang options
GO       ?= go
BINDATA  ?= bindata
MOCKGEN  ?= mockgen
PKG	 := github.com/nomkhonwaan/myblog
TAGS     :=
LDFLAGS  :=
GOFLAGS  :=
BINDIR   := $(CURDIR)/bin

# Node.js options
NPM	 ?= npm
NPX	 ?= npx
NG	 ?= $(NPX) ng

# Docker options
DOCKER   := docker

.PHONY: all
all: clean install build

.PHONY: install
install:
	$(GO) mod download

.PHONY: clean
clean:
	@rm -rf $(BINDIR)

.PHONY: mockgen
mockgen:
	$(MOCKGEN) -package blog -destination ./pkg/blog/category_mock.go github.com/nomkhonwaan/myblog/pkg/blog CategoryRepository
	$(MOCKGEN) -package blog -destination ./pkg/blog/post_mock.go github.com/nomkhonwaan/myblog/pkg/blog PostRepository
	$(MOCKGEN) -package blog -destination ./pkg/blog/tag_mock.go github.com/nomkhonwaan/myblog/pkg/blog TagRepository
	$(MOCKGEN) -package blog -destination ./pkg/blog/service_mock.go github.com/nomkhonwaan/myblog/pkg/blog Service
	$(MOCKGEN) -package mongo -destination ./pkg/mongo/collection_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Collection
	$(MOCKGEN) -package mongo -destination ./pkg/mongo/cursor_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Cursor
	$(MOCKGEN) -package mongo -destination ./pkg/mongo/single_result_mock.go github.com/nomkhonwaan/myblog/pkg/mongo SingleResult
	
.PHONY: test
test:
	$(GO) test -v ./... -race -coverprofile=coverage.txt -covermode=atomic
	
.PHONY: bindata
bindata:
	$(BINDATA) -o ./pkg/data/data.go ./data/...

.PHONY: build
build:
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '-X main.version=$(VERSION) -X main.revision=$(REVISION) $(LDFLAGS)' $(PKG)/cmd/myblog

.PHONY: build-web
build-web:
	cd web && $(NG) build --prod

.PHONY: docker-build
docker-build:
	docker build -f build/package/Dockerfile -t nomkhonwaan/myblog:latest .

.PHONY: deploy-web
deploy-web:
	cd web && \
	$(NPX) firebase use www-nomkhonwaan-com --token=${FIREBASE_TOKEN} && \
	$(NPX) firebase deploy --token=${FIREBASE_TOKEN}
