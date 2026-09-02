.PHONY: run test fmt

run:
	go run ./cmd/prompt-registry

test:
	go test ./...

fmt:
	gofmt -w cmd internal
