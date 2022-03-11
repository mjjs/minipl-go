test:
	cd src/pkg && go test ./... -timeout 2s
	cd src/cmd/minipl-go && go test ./... -timeout 2s

coverage:
	cd src/pkg && go test ./... -cover -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out

build:
	cd ./src/cmd/minipl-go && go build -o ../../../minipl-go

.PHONY: test coverage
