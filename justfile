check: test run-examples lint

# Run all tests including examples
test: generate
	go test ./...

run-examples:
	go run ./examples/tour
	go run ./examples/codegen
	go run ./examples/parse
	go run ./examples/parse-grpc

# Run benchmarks. May take a long time, use bench-short to avoid.
bench: generate
	go test ./... -bench=. -benchmem

# Run short benchmarks.
bench-short: generate
	go test ./... -bench=. -benchmem -short

# Generate test coverage
coverage:
	go test -coverprofile=coverage.html ./...
	go tool cover -html=coverage.html

# Generate code
generate: install-schemagen
	go generate ./...

# Install schemagen
install-schemagen:
	go install ./cmd/schemagen

# Update ISO datasets
update-iso: && generate
	curl -f -L -o ./internal/iso/countries.csv https://raw.githubusercontent.com/datasets/country-codes/refs/heads/main/data/country-codes.csv

	curl -f -L -o ./internal/iso/currencies.csv https://raw.githubusercontent.com/datasets/currency-codes/refs/heads/main/data/codes-all.csv

	curl -f -L -o ./internal/iso/languages.csv https://raw.githubusercontent.com/datasets/language-codes/refs/heads/main/data/language-codes-3b2.csv

lint:
	golangci-lint run --tests=false
