default: vet test

test:
	@go test -cover ./...

vet:
	@go vet -v ./...

.PHONY: test vet
