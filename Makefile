test:
	cd src/pkg && go test ./... -timeout 2s

coverage:
	cd src/pkg && go test ./... -cover -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out

.PHONY: test coverage
