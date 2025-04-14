test: generate
	go test ./...
	go run ./examples/tour
	go run ./examples/codegen

coverage:
	go test -coverprofile=coverage.html ./...
	go tool cover -html=coverage.html

generate:
	go generate ./...
