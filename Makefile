VERSION_PKG := github.com/ergoapi/util/version
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git branch -r --contains | head -1 | sed -E -e "s%(HEAD ->|origin|upstream)/?%%g" | xargs)
GIT_VERSION := $(shell git describe --always --tags --abbrev=14 $(GIT_COMMIT)^{commit})
GIT_TREE_STATE := $(shell s=`git status --porcelain 2>/dev/null`; if [ -z "$$s" ]; then echo "clean"; else echo "dirty"; fi)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := "-w \
	-X $(VERSION_PKG).gitVersion=$(GIT_VERSION) \
	-X $(VERSION_PKG).gitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PKG).gitBranch=$(GIT_BRANCH) \
	-X $(VERSION_PKG).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PKG).gitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PKG).gitMajor=0 \
	-X $(VERSION_PKG).gitMinor=0"

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## fmt code
	gofmt -s -w .
	goimports -w .
	@echo gofmt -l
	@OUTPUT=`gofmt -l . 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "gofmt must be run on the following files:"; \
        echo "$$OUTPUT"; \
        exit 1; \
    fi

lint: ## lint code
	@echo golangci-lint run --skip-files \".*test.go\" -v ./...
	@OUTPUT=`command -v golangci-lint >/dev/null 2>&1 && golangci-lint run --skip-files ".*test.go"  -v ./... 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "go lint errors:"; \
		echo "$$OUTPUT"; \
	fi

default: fmt lint ## fmt code

demo:
	cd example && go build -ldflags $(LDFLAGS)
