# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

PACKAGE = github.com/basgys/goxml2json

.PHONY: test386 test test-race fmt lint vet test-cover-html help
.DEFAULT_GOAL := help

test386: ## Run tests in 32-bit mode
	GOARCH=386 govendor test +local

test: ## Run tests
	govendor test +local

test-race: ## Run tests with race detector
	govendor test -race +local

fmt: ## Run gofmt linter
	@for d in `go list -no-status +local | sed 's/github.com.gohugoio.hugo/./'` ; do \
		if [ "`gofmt -l $$d/*.go | tee /dev/stderr`" ]; then \
			echo "^ improperly formatted go files" && echo && exit 1; \
		fi \
	done

lint: ## Run golint linter
	if [ "`golint . | tee /dev/stderr`" ]; then \
		echo "^ golint errors!" && echo && exit 1; \
	fi \

vet: ## Run go vet linter	
	@if [ "`go vet . | tee /dev/stderr`" ]; then \
		echo "^ go vet errors!" && echo && exit 1; \
	fi

test-cover-html: ## Generate test coverage report
	echo "mode: count" > coverage-all.out
	go test -coverprofile=coverage.out -covermode=count .
	tail -n +2 coverage.out >> coverage-all.out;
	go tool cover -html=coverage-all.out

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'