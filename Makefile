.PHONY: *

lint: lint-tp

# first-party linters
lint-fp:
	@echo ">>> running first-party linters ..."
	@go fmt ./...
	@go vet ./...
	@echo ">>> done"

# third-party linters
lint-tp:
	@echo ">>> running third-party linters ..."
	@go tool gofumpt -w ./...
	@go tool errcheck ./...
	@go tool staticcheck ./...
	@echo ">>> done"

golangci-lint:
	@echo ">>> running golangci-lint ..."
	@golangci-lint run ./...
	@echo ">>> done"
