default: test

test:
	@go test -v -cover -tags="integration" ./...

.PHONY: test
