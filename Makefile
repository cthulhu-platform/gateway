.PHONY: dev build run test lint clean stop

dev:
	go mod tidy
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing air for hot reloading..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@if ! command -v jq >/dev/null 2>&1; then \
		echo "Error: jq is required but not installed. Please install jq first."; \
		exit 1; \
	fi
	@echo "Starting development server with log tailing..."
	@echo "Logs will be written to /tmp/latest_run.log and displayed here with jq formatting"
	@echo "Press Ctrl+C to stop or run 'make stop'"
	@./scripts/dev.sh

build:
	@mkdir -p tmp/bin
	go build -o tmp/bin/api ./cmd/api

run:
	go run ./cmd/api

lint:
	golangci-lint run

stop:
	@echo "Stopping development server..."
	@./scripts/stop.sh

clean:
	rm -rf tmp/

test:
	go test ./... -v
