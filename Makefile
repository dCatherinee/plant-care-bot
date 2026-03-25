ifneq (,$(wildcard .env))
include .env
export
endif

GOBIN:=$(shell go env GOBIN)
GOPATH:=$(shell go env GOPATH)

ifeq ($(strip $(GOBIN)),)
GO_BIN_DIR:=$(GOPATH)/bin
else
GO_BIN_DIR:=$(GOBIN)
endif

GOOSE:=$(shell command -v goose 2>/dev/null)
ifeq ($(strip $(GOOSE)),)
GOOSE:=$(GO_BIN_DIR)/goose
endif

MIGRATIONS_DIR=./migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: test lint fmt vet migrate-up migrate-down

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint: fmt vet

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

test-integration:
	set -a; source .env; set +a; go test -tags=integration ./internal/storage/postgres -v

