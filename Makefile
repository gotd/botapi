test:
	@./go.test.sh
.PHONY: test

coverage:
	@./go.coverage.sh
.PHONY: coverage

lint:
	golangci-lint run ./...
.PHONY: lint

tidy:
	go mod tidy
.PHONY: tidy
