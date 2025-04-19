.PHONY: *

build:
	@echo ">>> building..."
	@GOOS=linux GOARCH=amd64 go build -o bootstrap
	@echo ">>> zipping..."
	@zip func.zip bootstrap
	@echo ">>> done"
