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
FIREBASE := $(NPX) firebase
NG	 := $(NPX) ng
WEBDIR   := $(CURDIR)/web

# Docker options
DOCKER   := docker

.PHONY: all
all: clean install build
	
.PHONY: install
install:
	$(GO) mod download

.PHONY: install-web
install-web:
	cd web && $(NPM) install --silent

.PHONY: clean
clean:
	@rm -rf $(BINDIR) && \
	@rm -rf $(CURDIR)/coverage.out && \
	@rm -rf $(CURDIR)/vendor && \
	@rm -rf $(WEBDIR)/dist && \
	@rm -rf $(WEBDIR)/.firebase && \
	@rm -rf $(WEBDIR)/node_modules

.PHONY: mockgen
mockgen:
	$(MOCKGEN) -destination ./pkg/auth/mock/client_mock.go net/http RoundTripper
	$(MOCKGEN) -destination ./pkg/blog/mock/category_mock.go github.com/nomkhonwaan/myblog/pkg/blog CategoryRepository
	$(MOCKGEN) -destination ./pkg/blog/mock/post_mock.go github.com/nomkhonwaan/myblog/pkg/blog PostRepository
	$(MOCKGEN) -destination ./pkg/blog/mock/tag_mock.go github.com/nomkhonwaan/myblog/pkg/blog TagRepository
	$(MOCKGEN) -destination ./pkg/graphql/mock/service_mock.go github.com/nomkhonwaan/myblog/pkg/graphql Service
	$(MOCKGEN) -destination ./pkg/mongo/mock/collection_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Collection
	$(MOCKGEN) -destination ./pkg/mongo/mock/cursor_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Cursor
	$(MOCKGEN) -destination ./pkg/mongo/mock/single_result_mock.go github.com/nomkhonwaan/myblog/pkg/mongo SingleResult
	$(MOCKGEN) -destination ./pkg/storage/mock/file_mock.go github.com/nomkhonwaan/myblog/pkg/storage FileRepository
	$(MOCKGEN) -destination ./pkg/storage/mock/storage_mock.go github.com/nomkhonwaan/myblog/pkg/storage Cache,Downloader,Uploader
	
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
	mv $(WEBDIR)/src/environments/environment.prod.ts $(WEBDIR)/src/environments/environment.prod.original.ts && \
	VERSION=$(VERSION) REVISION=$(REVISION) envsubst < $(WEBDIR)/src/environments/environment.prod.original.ts > $(WEBDIR)/src/environments/environment.prod.ts && \
	cd $(WEBDIR) && \
	$(NG) build --prod && \
	cd $(CURDIR) && \
	mv $(WEBDIR)/src/environments/environment.prod.original.ts $(WEBDIR)/src/environments/environment.prod.ts

.PHONY: build-docker
build-docker:
	$(DOCKER) build --file build/package/Dockerfile --tag nomkhonwaan/myblog:latest .

.PHONY: build-docker-all-in-one
build-docker-all-in-one:
	$(DOCKER) build --build-arg NPM_AUTH_TOKEN=${NPM_AUTH_TOKEN} --file build/package/Dockerfile.all-in-one --tag nomkhonwaan/myblog-all-in-one:latest .

.PHONY: deploy-web-firebase
deploy-web:
	cd $(WEBDIR) && \
	$(FIREBASE) use www-nomkhonwaan-com --token=${FIREBASE_TOKEN} && \
	$(FIREBASE) deploy --token=${FIREBASE_TOKEN}
