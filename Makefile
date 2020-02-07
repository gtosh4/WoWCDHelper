VERSION = $(shell echo $$(git describe --tags --always --abbrev=0 | sed -e 's/^v//' -e 's/\+.*//')+$$(git log -1 --pretty=format:%h)$$([ -z "$$(git status -s)" ] || echo ".d"))
TAG = $(shell echo $(VERSION) | sed 's/\+/-/g')

GO_FILES    ?= $(shell find . -name '*.go' -not -path './vendor/*')
GO_PACKAGES ?= $(shell go list ./...)
GOFLAGS     ?= -mod=vendor

FRONTEND_FILES ?= $(shell find . -path './frontend/src/*')

PB_FILES ?= $(shell find . -name '*.proto')

GOGO_REPLACE := Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types
PROTOC ?= protoc -I=. -I=../../../vendor/github.com/gogo/protobuf --gogofast_out=$(GOGO_REPLACE):.

.PHONY: all
all: bin/wcdh

bin/wcdh: proto fmt test $(GO_FILES) bin/frontend/dist vendor/modules.txt
	go build  -v -o $@ ./cmd/wcdh

frontend/node_modules: frontend/package.json
	cd frontend && npm install

frontend/dist/index.html: frontend/node_modules $(FRONTEND_FILES)
	cd frontend && npm run-script build

bin/frontend/dist: frontend/dist
	mkdir -p bin/frontend/
	ln -s ../../frontend/dist $@

vendor/modules.txt: go.mod go.sum
	go mod tidy
	go mod vendor

.PHONY: proto
proto: $(PB_FILES)
	cd pkg/warcraftlogs/events && $(PROTOC) *.proto
	cd pkg/warcraftlogs/fight && $(PROTOC) *.proto

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
