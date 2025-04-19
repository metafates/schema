test: generate
	go test ./...
	go run ./examples/tour
	go run ./examples/codegen

bench: generate
	go test ./... -bench=. -benchmem

coverage:
	go test -coverprofile=coverage.html ./...
	go tool cover -html=coverage.html

generate: install-schemagen
	go generate ./...

install-schemagen:
	go install ./cmd/schemagen
