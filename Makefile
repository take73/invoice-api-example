VERSION = 0.0.1
GO_PACKAGES = $(shell go list ./... | grep -v '/vendor/')

.PHONY: lint
lint:
	@go vet $(GO_PACKAGES)
	@staticcheck $(GO_PACKAGES)

.PHONY: wire
wire:
	cd cmd/server && wire

.PHONY: go run
run-local:
	go run cmd/server/main.go

.PHONY: migrate up
migrate-up:
	migrate -database "${DB_CONNECT_STRING_TEST}" -path db/migrations up

.PHONY: migrate down
migrate-down:
	migrate -database "${DB_CONNECT_STRING_TEST}" -path db/migrations down 1

.PHONY: test
test:
	go test -v -race $(GO_PACKAGES)
