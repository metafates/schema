test: generate
	go test ./...
	go run ./_examples/tour/main.go
	go run ./_examples/codegen/main.go

coverage:
	go test -coverprofile=coverage.html ./...
	go tool cover -html=coverage.html

generate:
	go generate ./...
