test:
	@echo "Testing..."
	@go test ./... -v
sec:
	@echo "Running Go Sec..."
	@gosec ./...
lint:
	@echo "Running Linter..."
	@staticcheck ./...

docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi
