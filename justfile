check: fmt test run-examples lint

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
	# country codes
	curl -f -L -o ./internal/iso/countries.csv https://raw.githubusercontent.com/datasets/country-codes/refs/heads/main/data/country-codes.csv

	# currencies
	curl -f -L -o ./internal/iso/currencies.csv https://raw.githubusercontent.com/datasets/currency-codes/refs/heads/main/data/codes-all.csv

	# languages
	curl -f -L -o ./internal/iso/languages.csv https://raw.githubusercontent.com/datasets/language-codes/refs/heads/main/data/language-codes-3b2.csv

# format source code
fmt:
	golangci-lint fmt

# lint source code
lint:
	golangci-lint run --tests=false

# Open documentation
doc:
	go run golang.org/x/pkgsite/cmd/pkgsite@latest -open
