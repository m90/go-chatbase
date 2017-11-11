default: vet test

test:
	@go test -cover ./...

vet:
	@go vet ./...

.PHONY: test vet
