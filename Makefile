default: vet test

test:
	@go test -v -cover -tags="integration" ./...

vet:
	@go vet ./...

.PHONY: test vet
