.PHONY: help
help: ## Print help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-lint-tools: ## Install lint tools
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/alexkohler/prealloc@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

DIR=./...
lint: ## Run static analysis
	go vet "$(DIR)"
	test -z "`gofmt -s -d .`"
	staticcheck "$(DIR)"
	prealloc -set_exit_status "$(DIR)"
	gosec "$(DIR)"

.PHONY: test
test: ## run test ex.) make test OPT="-run TestXXX"
	go test -v ./... "$(OPT)"

test-coverage: ## Run test with coverage
	$(MAKE) test OPT="-coverprofile=coverage.out"
	go tool cover -html=coverage.out
