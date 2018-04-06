setup: ## Install all the build and lint dependencies
	go get -u github.com/golang/dep/...
	go get github.com/goreleaser/goreleaser
	dep ensure

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

ci: snapshot

build: clean ## Build a version
	go build -o dist/download-rds-logs

install: ## Install to $GOPATH/src
	go install

run: ## Run main.g
	go run main.go

clean: ## Clean any builds
	rm -rf dist

release: clean ## Prepare a release
	goreleaser

snapshot: clean ## Prepare a snapshot
	goreleaser --snapshot


# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
