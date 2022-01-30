test:
	go test ./...

coverage:
	go test ./... -cover -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out

.PHONY: test coverage
