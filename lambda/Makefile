.PHONY: *

run:
	@echo ">>> running..."
	@go run main.go
	@echo ">>> done"

build:
	@echo ">>> building..."
	@GOOS=linux GOARCH=amd64 go build -o bootstrap
	@echo ">>> zipping..."
	@zip func.zip bootstrap
	@echo ">>> done"
