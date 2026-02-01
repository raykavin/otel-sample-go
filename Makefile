.PHONY: help test run

help:
	@echo "Available targets:"
	@echo "  run             Run the Go application"
	@echo "  test            Run Go app and execute test.sh (Vegeta attack)"

run:
	@echo "Starting Go application..."
	@go run .

test:
	@echo "Running Go application..."
	@go run . & \
	APP_PID=$$!; \
	sleep 2; \
	echo "Running tests..."; \
	scripts/sh/test.sh; \
	echo "Stopping Go application..."; \
	kill $$APP_PID
