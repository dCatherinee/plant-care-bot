.PHONY: test lint fmt vet

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint: fmt vet