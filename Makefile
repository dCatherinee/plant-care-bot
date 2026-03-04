.PHONY: test lint fmt vet

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet