default: test

test:
ifeq ($(TRAVIS_PULL_REQUEST),true)
	# do not run integration tests for pull requests as they would
	# require the secret chatbase api key to be present
	@go test -v -cover ./...
else
	@go test -v -cover -tags="integration" ./...
endif

.PHONY: test
