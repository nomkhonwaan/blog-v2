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
	$(MOCKGEN) -destination ./pkg/auth/mock/client_mock.go net/http RoundTripper
	$(MOCKGEN) -destination ./pkg/blog/mock/category_mock.go github.com/nomkhonwaan/myblog/pkg/blog CategoryRepository
	$(MOCKGEN) -destination ./pkg/blog/mock/post_mock.go github.com/nomkhonwaan/myblog/pkg/blog PostRepository
	$(MOCKGEN) -destination ./pkg/blog/mock/tag_mock.go github.com/nomkhonwaan/myblog/pkg/blog TagRepository
	$(MOCKGEN) -destination ./pkg/blog/mock/service_mock.go github.com/nomkhonwaan/myblog/pkg/blog Service
	$(MOCKGEN) -destination ./pkg/mongo/mock/collection_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Collection
	$(MOCKGEN) -destination ./pkg/mongo/mock/cursor_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Cursor
	$(MOCKGEN) -destination ./pkg/mongo/mock/single_result_mock.go github.com/nomkhonwaan/myblog/pkg/mongo SingleResult
	
.PHONY: test
test:
	$(GO) test -v ./... -race -coverprofile=coverage.out -covermode=atomic
	
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
