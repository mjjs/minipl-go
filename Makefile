test:
	cd src && go test ./...

coverage:
	cd src && go test ./... -cover -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out

.PHONY: test coverage
