test:
	go test ./..

lint: fmt vet

fmt:
	go fmt ./..

vet:
	go vet ./..