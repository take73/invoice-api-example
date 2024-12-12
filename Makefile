VERSION = 0.0.1
GO_PACKAGES = $(shell go list ./... | grep -v '/vendor/')

.PHONY: lint
lint:
	@go vet $(GO_PACKAGES)
	@staticcheck $(GO_PACKAGES)

.PHONY: go run
run-local:
	@set -a && source .env && go run cmd/server/main.go

.PHONY: migrate up
migrate-up:
	@set -a && source .env && migrate -database "$${DB_CONNECT_STRING_TEST}" -path db/migrations up

.PHONY: migrate down
migrate-down:
	@set -a && source .env && migrate -database "$${DB_CONNECT_STRING_TEST}" -path db/migrations down 1

.PHONY: test
test:
	@set -a && source .env && go test -v -race $(GO_PACKAGES)

.PHONY: export
export:
	cd tmp/gpt-repository-loader && python gpt_repository_loader.py ../../internal -o example.txt