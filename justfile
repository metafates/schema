# Run all tests including examples
test: generate
	go test ./...
	go run ./examples/tour
	go run ./examples/codegen
	go run ./examples/parse

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
