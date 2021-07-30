VERSION = $(shell echo $$(git describe --tags --always --abbrev=0 | sed -e 's/^v//' -e 's/\+.*//')+$$(git log -1 --pretty=format:%h)$$([ -z "$$(git status -s)" ] || echo ".d"))
TAG = $(shell echo $(VERSION) | sed 's/\+/-/g')

GO_FILES    ?= $(shell find . -name '*.go' -not -path './vendor/*')
GO_PACKAGES ?= $(shell go list ./...)
GOFLAGS     ?= -mod=vendor

FRONTEND_FILES ?= $(shell find . -path './frontend/src/*')

.PHONY: all
all: bin/wcdh

bin/wcdh: fmt test $(GO_FILES) bin/frontend/dist go.sum
	go build  -v -o $@ ./cmd/wcdh

frontend/node_modules: frontend/package.json
	cd frontend && npm install

frontend/dist/index.html: frontend/node_modules $(FRONTEND_FILES)
	cd frontend && npm run-script build

bin/frontend/dist: frontend/dist
	mkdir -p bin/frontend/
	ln -s ../../frontend/dist $@

go.sum: go.mod
	go mod tidy

.PHONY: test
test:
	@go test -race $(GO_PACKAGES)

.PHONY: fmt
fmt: $(GO_FILES)
	@gofmt -w $(GO_FILES)

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: image
image: test
	docker build -t wowcds:$(TAG) -f build/package/Dockerfile .
	docker tag wowcds:$(TAG) wowcds:latest
