default: vet test

test:
	@go test -cover -tags="integration" ./...

vet:
	@go vet ./...

.PHONY: test vet
